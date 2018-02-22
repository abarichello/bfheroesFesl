package network

import (
	"bytes"
	"errors"
	"io"
	"net"

	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

func (client *Client) Answer(Pkt *codec.Pkt) error {
	if !client.IsActive {
		logrus.Printf("%s: Trying to write to inactive Client.\n%v", client.name, Pkt.Type)
		return errors.New("Client NOT active.Can't send message")
	}

	return Answer(Pkt, func(buf *bytes.Buffer) error {
		_, err := io.Copy(client.conn, buf)
		return err
	})
}

func (socket *SocketUDP) Answer(Pkt *codec.Pkt, addr *net.UDPAddr) error {
	return Answer(Pkt, func(buf *bytes.Buffer) error {
		_, err := socket.listen.WriteToUDP(buf.Bytes(), addr)
		return err
	})
}

func Answer(Pkt *codec.Pkt, writer func(*bytes.Buffer) error) error {
	logger := logrus.WithFields(logrus.Fields{"type": Pkt.Type, "HEX": Pkt.Send})

	encoder := codec.NewEncoder()
	buf, err := encoder.EncodePkt(Pkt)
	if err != nil {
		logger.WithError(err).Error("Cannot encode Pkt")
		return nil
	}

	err = writer(buf)
	if err != nil {
		logger.WithError(err).Error("Cannot write Pkt")
		return nil
	}
	return nil
}
