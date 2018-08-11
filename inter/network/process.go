package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strings"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type eventReadFesl func(outCommand *ProcessFESL, ContentType string)

func (client *Client) readFESL(data []byte) []byte {
	return readFesl(data, func(cmd *ProcessFESL, ContentType string) {
		client.eventChan <- ClientEvent{Name: "command." + ContentType, Data: cmd}
		client.eventChan <- ClientEvent{Name: "command", Data: cmd}
	})
}

func (client *Client) readFESLTLS(data []byte) []byte {
	return readFesl(data, func(cmd *ProcessFESL, ContentType string) {
		client.eventChan <- ClientEvent{Name: "command." + cmd.Msg["TXN"], Data: cmd}
		client.eventChan <- ClientEvent{Name: "command", Data: cmd}
	})
}


//new
// func (client *Client) readFESL(data []byte) {
// 	cmds, err := codec.ParseCommands(data)
// 	if err != nil {
// 		logrus.
// 			WithError(err).
// 			WithField("packet", string(data)).
// 			Error("Cannot parse commands")
// 		return
// 	}
// 	for _, cmd := range cmds {
// 		client.receiver <- ClientEvent{
// 			Name: "command." + cmd.Query,
// 			Data: cmd,
// 		}
// 	}
// }


// func (socket *SocketUDP) readFESL(data []byte, addr *net.UDPAddr) {
// 	p := bytes.NewBuffer(data)
// 	var HEX uint32
// 	var length uint32

// 	ContentType := string(data[:4])
// 	p.Next(4)

// 	binary.Read(p, binary.BigEndian, &HEX)
// 	binary.Read(p, binary.BigEndian, &length)

// 	ContentRaw := data[12:]
// 	Content := codec.DecodeFESL(ContentRaw)

// 	out := &ProcessFESL{
// 		Query: ContentType,
// 		HEX:   HEX,
// 		Msg:   Content,
// 	}

// 	socket.EventChan <- SocketUDPEvent{Name: "command." + ContentType, Addr: addr, Data: out}
// 	socket.EventChan <- SocketUDPEvent{Name: "command", Addr: addr, Data: out}
// }

func (client *Client) readTLSPacket(data []byte) {
	cmds, err := codec.ParseCommands(data)
	if err != nil {
		var extract string
		if len(data) > 128 {
			extract = string(data[:128])
		} else {
			extract = string(data)
		}

		logrus.
			WithError(err).
			WithField("extract", extract).
			Error("Cannot parse commands (TLS)")
		return
	}
	for _, cmd := range cmds {
		client.receiver <- ClientEvent{
			Name: "command." + cmd.Message["TXN"],
			Data: cmd,
		}
	}
}

type RawPacket struct {
	// Query first 4 bytes
	// i.e. "fsys", "acct" "CONN", "UPLA"
	Query []byte

	// Broadcast, next 4 bytes
	// TOOD: Use enumerator
	Broadcast []byte

	// Length, next 4 bytes
	Length []byte

	// Payload
	Payload []byte
}


func (socket *SocketUDP) readFESL(data []byte, addr *net.UDPAddr) {
	p := bytes.NewBuffer(data)
	var payloadID uint32
	var payloadLen uint32

	payloadType := string(data[:4])
	p.Next(4)

	binary.Read(p, binary.BigEndian, &payloadID)
	binary.Read(p, binary.BigEndian, &payloadLen)

	payload := codec.DecodeFESL(data[12:])

	socket.EventChan <- SocketUDPEvent{
		Name: payloadType,
		Addr: addr,
		Data: &codec.Command{
			Query:     payloadType,
			PayloadID: payloadID,
			Message:   payload,
		},
	}
}



func readFesl(data []byte, fireEvent eventReadFesl) []byte {
	p := bytes.NewBuffer(data)
	i := 0
	
	var ContentRaw []byte
	for {
		// Create a copy at this point in case we have to abort later
		// And send back the Packet to get the rest
		curData := p

		var HEX uint32
		var length uint32

		ContentTypeRaw := make([]byte, 4)
		_, err := p.Read(ContentTypeRaw)
		if err != nil {
			return nil
		}

		ContentType := string(ContentTypeRaw)

		binary.Read(p, binary.BigEndian, &HEX)


		if p.Len() < 4 {
			return nil
		}


		binary.Read(p, binary.BigEndian, &length)

		if (length - 12) > uint32(len(p.Bytes())) {
			logrus.Println("Packet not fully read")
			return curData.Bytes()
		}

		ContentRaw = make([]byte, (length - 12))
		p.Read(ContentRaw)

		Content := codec.DecodeFESL(ContentRaw) //hex to asci

		out := &ProcessFESL{
			Query: ContentType,
			HEX:   HEX, //ContentID like 0xc000000d
			Msg:   Content,
		}
		fireEvent(out, ContentType)

		i++
	}

	return nil
}

type ProcessFESL struct {
	Msg   map[string]string
	Query string
	HEX   uint32
}

// processCommand turns gamespy's command string to the
// command struct
func processCommand(msg string) (*ProcessFESL, error) {
	
	
	
	outCommand := new(ProcessFESL) // Command not a CommandFESL
	outCommand.Msg = make(map[string]string)
	data := strings.Split(msg, `\`)

	if len(data) <= 0 {
		logrus.Errorln("Command Msg invalid")
		return nil, errors.New("Command Msg invalid")
	}

	if len(data) == 1 {
		outCommand.Msg["__query"] = data[0]
		outCommand.Query = data[0]
		return outCommand, nil
	}

	outCommand.Query = data[1]
	outCommand.Msg["__query"] = data[1]
	for i := 1; i < len(data)-1; i = i + 2 {
		outCommand.Msg[strings.ToLower(data[i])] = data[i+1]
	}

	return outCommand, nil
}

func (client *Client) processCommand(command string) {
	gsPacket, err := processCommand(command)
	if err != nil {
		logrus.Errorf("%s: Error processing command %s.\n%v", client.name, command, err)
		client.eventChan <- client.FireError(err)
		return
	}

	client.eventChan <- ClientEvent{Name: "command." + gsPacket.Query, Data: gsPacket}
	client.eventChan <- ClientEvent{Name: "command", Data: gsPacket}
}

func (socket *SocketUDP) processCommand(command string, addr *net.UDPAddr) {
	gsPacket, err := processCommand(command)
	if err != nil {
		logrus.Errorf("%s: Error processing command %s.\n%v", socket.name, command, err)
		socket.EventChan <- SocketUDPEvent{Name: "error", Addr: addr, Data: err}
		return
	}

	socket.EventChan <- SocketUDPEvent{Name: "command." + gsPacket.Query, Addr: addr, Data: gsPacket}
	socket.EventChan <- SocketUDPEvent{Name: "command", Addr: addr, Data: gsPacket}
}
