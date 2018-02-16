package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansStart struct {
	Txn    string `fesl:"TXN"`
	ID     string `fesl:"id.id"`
	idpart string `fesl:"id.partition"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientCommand) {
	event.Client.WriteEncode(&codec.Packet{
		Payload: ansStart{
			Txn:    "Start",
			ID:     "1",
			idpart: event.Command.Message["partition.partition"],
		},
		Type: pnow,
	})

	fm.Status(event)
	logrus.Println("=START=")
}
