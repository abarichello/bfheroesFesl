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

type eventReadFesl func(outCommand *CommandFESL, payloadType string)

func (client *Client) readFESL(data []byte) []byte {
	return readFesl(data, func(cmd *CommandFESL, payloadType string) {
		client.eventChan <- ClientEvent{Name: "command." + payloadType, Data: cmd}
		client.eventChan <- ClientEvent{Name: "command", Data: cmd}
	})
}

func (client *Client) readFESLTLS(data []byte) []byte {
	return readFesl(data, func(cmd *CommandFESL, payloadType string) {
		client.eventChan <- ClientEvent{Name: "command." + cmd.Msg["TXN"], Data: cmd}
		client.eventChan <- ClientEvent{Name: "command", Data: cmd}
	})
}

func (socket *SocketUDP) readFESL(data []byte, addr *net.UDPAddr) {
	p := bytes.NewBuffer(data)
	var payloadID uint32
	var payloadLen uint32

	payloadType := string(data[:4])
	p.Next(4)

	binary.Read(p, binary.BigEndian, &payloadID)
	binary.Read(p, binary.BigEndian, &payloadLen)

	payloadRaw := data[12:]
	payload := codec.DecodeFESL(payloadRaw)

	out := &CommandFESL{
		Query:     payloadType,
		PayloadID: payloadID,
		Msg:       payload,
	}

	socket.EventChan <- SocketUDPEvent{Name: "command." + payloadType, Addr: addr, Data: out}
	socket.EventChan <- SocketUDPEvent{Name: "command", Addr: addr, Data: out}
}

func readFesl(data []byte, fireEvent eventReadFesl) []byte {
	p := bytes.NewBuffer(data)
	i := 0
	var payloadRaw []byte
	for {
		// Create a copy at this point in case we have to abort later
		// And send back the packet to get the rest
		curData := p

		var payloadID uint32
		var payloadLen uint32

		payloadTypeRaw := make([]byte, 4)
		_, err := p.Read(payloadTypeRaw)
		if err != nil {
			return nil
		}

		payloadType := string(payloadTypeRaw)

		binary.Read(p, binary.BigEndian, &payloadID)

		if p.Len() < 4 {
			return nil
		}

		binary.Read(p, binary.BigEndian, &payloadLen)

		if (payloadLen - 12) > uint32(len(p.Bytes())) {
			logrus.Println("Packet not fully read")
			return curData.Bytes()
		}

		payloadRaw = make([]byte, (payloadLen - 12))
		p.Read(payloadRaw)

		payload := codec.DecodeFESL(payloadRaw)

		out := &CommandFESL{
			Query:     payloadType,
			PayloadID: payloadID,
			Msg:       payload,
		}
		fireEvent(out, payloadType)

		i++
	}

	return nil
}

type CommandFESL struct {
	Msg       map[string]string
	Query     string
	PayloadID uint32
}

// processCommand turns gamespy's command string to the
// command struct
func processCommand(msg string) (*CommandFESL, error) {
	outCommand := new(CommandFESL) // Command not a CommandFESL
	outCommand.Msg = make(map[string]string)
	data := strings.Split(msg, `\`)

	// TODO:
	// Should maybe return an emtpy Command struct instead
	if len(data) < 1 {
		logrus.Errorln("Command Msg invalid")
		return nil, errors.New("Command Msg invalid")
	}

	// TODO:
	// Check if that makes any sense? Kinda just translated from the js-code
	//		if (data.length < 2) { return out; }
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
