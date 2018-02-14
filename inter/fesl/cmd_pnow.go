package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/matchmaking"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

const (
	pnow = "pnow"
	// pnowCancel    = "Cancel"
	pnowStart  = "Start"
	pnowStatus = "Status"
)

type ansStatus struct {
	Taxon        string                 `fesl:"TXN"`
	ID           stPartition            `fesl:"id"`
	SessionState string                 `fesl:"sessionState"`
	Properties   map[string]interface{} `fesl:"props"`
}

//stPartition=statusPartition
type stPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

type statusGame struct {
	LobbyID int    `fesl:"lit"`
	Fit     int    `fesl:"fit"`
	GameID  string `fesl:"gid"`
}

// Status Call to overall status(BEFORE START)
func (fm *FeslManager) Status(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("C Left")
		return
	}

	gameID := matchmaking.FindGIDs()

	ans := ansStatus{
		Taxon:        pnowStatus,
		ID:           stPartition{1, event.Command.Message["partition.partition"]},
		SessionState: "COMPLETE",
		Properties: map[string]interface{}{
			"resultType": "JOIN",
			"games": []statusGame{
				{
					LobbyID: 1,
					Fit:     1001,
					GameID:  gameID,
				},
			},
		},
	}

	event.Client.WriteEncode(&codec.Packet{
		Payload: ans,
		Step:    0x80000000,
		Type:    pnow,
	})
}

type ansStart struct {
	Taxon string      `fesl:"TXN"`
	ID    stPartition `fesl:"id"`
}

//Start TODO SYNC W/Discord & HWID
func (fm *FeslManager) Start(event network.EventClientCommand) {

	event.Client.WriteEncode(&codec.Packet{
		Payload: ansStart{
			Taxon: pnowStart,
			ID:    stPartition{1, event.Command.Message["partition.partition"]},
		},
		Step: event.Command.PayloadID,
		Type: pnow,
	})

	fm.Status(event)
}
