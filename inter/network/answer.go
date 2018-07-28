package network

import (
	"bytes"
	//"errors"
	//"io"
	"net"

	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

func (socket *SocketUDP) Answer(Packet *codec.Packet, addr *net.UDPAddr) error {
	return AnswerUDP(Packet, func(buf *bytes.Buffer) error {
		_, err := socket.listen.WriteToUDP(buf.Bytes(), addr)
		return err
	})
}

func AnswerUDP(Packet *codec.Packet, writer func(*bytes.Buffer) error) error {
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

func (client *Client) SendPacket(pkt []byte) error {
	_, err := client.conn.Write(pkt)
	if err != nil {
		logrus.
			WithError(err).
			Warn("Cannot send encoded packet")
		return err
	}

	logrus.
		WithField("packet", string(pkt)).
		Print("client.SendPacket")

	return nil
}

func (client *Client) Answer(Packet *codec.Packet) error {
	if !client.IsActive {
		logrus.Println("Trying to write to inactive Client.\n%v", Packet.Message)
	}

	encoder := codec.NewEncoder()
	buf, err := encoder.EncodePacket(Packet)
	if err != nil {
		logrus.
			WithError(err).
			WithField("Message", Packet.Message).
			Error("Cannot encode packet")
		return err
	}

	return client.SendPacket(buf.Bytes())
}
