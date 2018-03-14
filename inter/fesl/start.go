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
	ID         stPartition `fesl:"id"`
	debugLevel string      `fesl:"debugLevel"`
	TXN        string      `fesl:"TXN"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")

	var Button string

	if event.Process.Msg["Start"] == "1" {
		Button = "Start"
	} else {
		Button = ""
	}

	event.Client.Answer(&codec.Pkt{
		Content: Start{
			TXN:        Button,
			debugLevel: "high",
			ID: stPartition{1,
				event.Process.Msg[partition]},
		},
		Send:    event.Process.HEX,
		Message: "pnow",
	})
	fm.Status(event)
}
