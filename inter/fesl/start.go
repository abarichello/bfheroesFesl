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
	TXN              string      `fesl:"TXN"`
	ID               stPartition `fesl:"id"`
	debugLevel       string      `fesl:"debugLevel"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")

	event.Client.Answer(&codec.Pkt{
		Content: Start{
			TXN:        "Start",
			debugLevel: "high",
			ID: stPartition{1,
				event.Process.Msg[partition]},
		},
		Send: 0xc000000d,
		Type: "pnow",
	})
	fm.Status(event)
}
