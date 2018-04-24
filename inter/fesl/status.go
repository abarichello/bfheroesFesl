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
	GID     string `fesl:"gid"` //gameID
}

type stPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

// Status comes after Start. tells info about desired server
func (fm *Fesl) Status(event network.EvProcess) {
	logrus.Println("=Status=")
	reply := event.Process.Msg
	gameID := mm.FindGIDs()

	games := []stGame{
		{
			LobbyID: 1,
			GID:     gameID,
			Fit:     1001,
		},
	}
	event.Client.Answer(&codec.Packet{
		Content: Status{
			TXN:   "Status",
			State: "COMPLETE",
			ID:    stPartition{1, reply[partition]},
			Props: map[string]interface{}{
				"resultType": "JOIN",
				"games":      games,
			},
		},
		Send:    0x80000000,
		Message: "pnow",
	})
}

type Cancel struct {
	TXN   string                 `fesl:"TXN"`
	ID    stPartition            `fesl:"id"`
	State string                 `fesl:"sessionState"`
	Props map[string]interface{} `fesl:"props"`
}

// Cancel - cancel pnow
func (fm *Fesl) Cancel(event network.EvProcess) {
	reply := event.Process.Msg

	event.Client.Answer(&codec.Packet{
		Content: Cancel{
			TXN:   "Cancel",
			State: "CANCELLED",
			ID:    stPartition{1, reply[partition]},
			Props: map[string]interface{}{
				"resultType": "CANCEL",
			},
		},
		Send:    event.Process.HEX,
		Message: event.Process.Query,
	})
}
