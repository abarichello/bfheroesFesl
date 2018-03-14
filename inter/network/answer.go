package network

import (
	"bytes"
	"errors"
	"io"
	"net"

	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

func (client *Client) Answer(Packet *codec.Packet) error {
	if !client.IsActive {
		logrus.Printf("%s: Trying to write to inactive Client.\n%v", client.name, Packet.Message)
		return errors.New("Client NOT active.Can't send message")
	}

	return Answer(Packet, func(buf *bytes.Buffer) error {
		_, err := io.Copy(client.conn, buf)
		return err
	})
}

func (socket *SocketUDP) Answer(Packet *codec.Packet, addr *net.UDPAddr) error {
	return Answer(Packet, func(buf *bytes.Buffer) error {
		_, err := socket.listen.WriteToUDP(buf.Bytes(), addr)
		return err
	})
}

func Answer(Packet *codec.Packet, writer func(*bytes.Buffer) error) error {
	logger := logrus.WithFields(logrus.Fields{"type": Packet.Message, "HEX": Packet.Send})

	encoder := codec.NewEncoder()
	buf, err := encoder.EncodePacket(Packet)
	if err != nil {
		logger.WithError(err).Error("Cannot encode Packet")
		return nil
	}

	err = writer(buf)
	if err != nil {
		logger.WithError(err).Error("Cannot write Packet")
		return nil
	}
	return nil
}
