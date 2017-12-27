package network

import (
	"bytes"
	"errors"
	"io"
	"net"

	"bitbucket.org/openheroes/backend/internal/network/codec"

	"github.com/sirupsen/logrus"
)

func (client *Client) WriteEncode(packet *codec.Packet) error {
	if !client.IsActive {
		logrus.Printf("%s: Trying to write to inactive Client.\n%v", client.name, packet.Type)
		return errors.New("Client NOT active.Can't send message")
	}

	return writeEncode(packet, func(buf *bytes.Buffer) error {
		_, err := io.Copy(client.conn, buf)
		return err
	})
}

func (socket *SocketUDP) WriteEncode(packet *codec.Packet, addr *net.UDPAddr) error {
	return writeEncode(packet, func(buf *bytes.Buffer) error {
		_, err := socket.listen.WriteToUDP(buf.Bytes(), addr)
		return err
	})
}

func writeEncode(packet *codec.Packet, writer func(*bytes.Buffer) error) error {
	logger := logrus.WithFields(logrus.Fields{"type": packet.Type, "payloadID": packet.Step})

	encoder := codec.NewEncoder()
	buf, err := encoder.EncodePacket(packet)
	if err != nil {
		logger.WithError(err).Error("Cannot encode packet")
		return nil
	}

	err = writer(buf)
	if err != nil {
		logger.WithError(err).Error("Cannot write packet")
		return nil
	}
	return nil
}
