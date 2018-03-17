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
	TXN        string                 `fesl:"TXN"`
	ID         stPartition            `fesl:"id"`
	State      string                 `fesl:"sessionState"`
	Props      map[string]interface{} `fesl:"props"`
}

type stGame struct {
	LobbyID int    `fesl:"lid"`
	Fit     int    `fesl:"fit"`
	GID     string `fesl:"gid"`
}

type stPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

// Status comes after Start. tells info about desired server
func (fm *FeslManager) Status(event network.EventClientProcess) {
	hex := event.Process.HEX
	reply := event.Process.Msg
	answer := event.Client.Answer
	logrus.Println("=Status=")
	gameID := mm.FindGIDs()

	//@TODO refactor this
	ans := Status{
		TXN: "Status",
		ID: stPartition{1,
			reply[partition]},
		State: "COMPLETE",
		Props: map[string]interface{}{
			"debugLevel": "low",
			"resultType": "JOIN",
			"games": []stGame{
				{
					LobbyID: 1,
					Fit:     1000,
					GID:     gameID,
				},
			},
		},
	}
	answer(&codec.Packet{
		Content: ans,
		Send:    hex,
		Message: "pnow",
	})
}
