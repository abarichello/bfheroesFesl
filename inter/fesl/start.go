package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type Start struct {
	TXN string      `fesl:"TXN"`
	ID  stPartition `fesl:"id"`
	dbgLevel string `fesl:"debugLevel"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")

	event.Client.Answer(&codec.Pkt{
		Content: Start{
			TXN: "Start",
			dbgLevel:  "high",
			ID: stPartition{1, event.Process.Msg[partition]},
		},
		Send: 0xc000000d,
		Type: "pnow",
	})
	fm.Status(event)
}
