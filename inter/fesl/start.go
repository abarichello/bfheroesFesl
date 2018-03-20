package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)
//TODO
// 'GetStatus'
// 'Update'
// 'Cancel'

type Start struct {
	ID  stPartition `fesl:"id"`
	TXN string      `fesl:"TXN"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")

	event.Client.Answer(&codec.Packet{
		Content: Start{
			TXN: "Start",
			ID: stPartition{1,
				event.Process.Msg[partition]},
		},
		Send:    event.Process.HEX,
		Message: "pnow",
	})
	fm.Status(event)
}

