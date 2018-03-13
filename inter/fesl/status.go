package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	partition = "partition.partition"
)

type Status struct {
	TXN   string                 `fesl:"TXN"`
	ID    stPartition            `fesl:"id"`
	State string                 `fesl:"sessionState"`
	Props map[string]interface{} `fesl:"props"`
}

type stPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

// Status comes after Start. tells info about desired server
func (fm *FeslManager) Status(event network.EventClientProcess) {
	logrus.Println("==Status==")
	gameID := mm.FindGIDs()

	//@TODO refactor this
	ans := Status{
		TXN: "Status",
		ID: stPartition{1,
			event.Process.Msg[partition]},
		State: "COMPLETE",
		Props: map[string]interface{}{
			"resultType":  "JOIN",
			"sessionType": "findServer",
			"games":       "1",
			"Fit":         "1000", //ELO
			"GID":         gameID,
		},
	}
	event.Client.Answer(&codec.Pkt{
		Content: ans,
		Send:    0x80000000,
		Type:    "pnow",
	})
}
