package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	Complete = "COMPLETE"
	pnow     = "pnow"
	//pnowCancel = "Cancel"
	pnowStart  = "Start"
	pnowStatus = "Status"
	partition  = "partition.partition"
	Join       = "JOIN"
)

type ansStatus struct {
	TXN        string                 `fesl:"TXN"`
	ID         stPartition            `fesl:"id"`
	State      string                 `fesl:"sessionState"`
	Properties map[string]interface{} `fesl:"props"`
}

type stPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

type stGame struct {
	LobbyID int    `fesl:"lid"`
	Fit     int    `fesl:"fit"` // ELO ?
	GAME    string `fesl:"gid"`
}

// Status pnow.Status command
func (fm *FeslManager) Status(event network.EventClientProcess) {
	logrus.Println("=Status=")
	//Infinite Search
	gameID := mm.FindGIDs()

	ans := ansStatus{
		TXN: pnowStatus,
		ID: stPartition{1,
			event.Process.Msg[partition]},
		State: Complete,
		Properties: map[string]interface{}{
			"resultType": Join,
			"games": []stGame{
				{
					LobbyID: 1,
					Fit:     1500,
					GAME:    gameID,
				},
			},
		},
	}
	event.Client.Answer(&codec.Pkt{
		Content: ans,
		Send:    0x80000000,
		Type:    pnow,
	})
}
