package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansStart struct {
	TXN string      `fesl:"TXN"`
	ID  stPartition `fesl:"id"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")
	event.Client.Answer(&codec.Pkt{
		Content: ansStart{
			TXN: pnowStart,
			ID: stPartition{
				1,
				event.Process.Msg[partition]},
		},
		Send: event.Process.HEX,
		Type: pnow,
	})
	fm.Status(event)
}
