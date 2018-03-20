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
	logrus.Println("=Status=")
	reply := event.Process.Msg
	answer := event.Client.Answer
	gameID := mm.FindGIDs()

	games := []stGame{
		{
			LobbyID:   1,
			GID:    gameID,
			Fit: 1001,
		},
	}

	payload := Status{
		TXN: "Status",
		ID: stPartition{1, reply[partition]},
		State: "COMPLETE",
		Props: map[string]interface{}{
			"resultType": "JOIN",
			"games": games,				
			},
		}

		answer(&codec.Packet{
		Content: payload,
		Send:    0x80000000,
		Message: "pnow",
	})
}

type Cancel struct {
	TXN    string      `fesl:"TXN"`
	ID     stPartition `fesl:"id"`
	State  string      `fesl:"sessionState"`
	Props map[string]interface{} `fesl:"props"`

}

// Cancel - cancel pnow
func (fm *FeslManager) Cancel(event network.EventClientProcess) {
	logrus.Println("==Cancel Q==")
	reply := event.Process.Msg

	event.Client.Answer(&codec.Packet{
		Content: Cancel{
			TXN: "Cancel",
			State: "CANCELLED",
			ID: stPartition{1, reply[partition]},
			Props: map[string]interface{}{
				"resultType": "CANCEL",
			},			
		},
		Send:    event.Process.HEX,
		Message: event.Process.Query,
	})
}

