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
	TXN            string      `fesl:"TXN"`
	ID             stPartition `fesl:"id"`
	enableEasyZone string      `fesl:"enableEasyZone"`
	debugLevel     string      `fesl:"debugLevel"`
	timeout        string      `fesl:"poolTimeout"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")
	ans := make(map[string]string)
	answer := event.Client.Answer
	ans["TXN"] = "Start"
	ans["id.id"] = "1"
	ans["id.partition"] = event.Process.Msg[partition]
	ans["poolTimeout"] = "60"
	ans["firewallType"] = "0"
	answer(&codec.Packet{
		Content: ans,
		Send:    0x80000000,
		Message: "pnow",
	})
	fm.Status(event)
}
