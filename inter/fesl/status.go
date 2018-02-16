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
	Txn          string                 `fesl:"TXN"`
	ID           stPartition            `fesl:"id"`
	SessionState string                 `fesl:"sessionState"`
	Properties   map[string]interface{} `fesl:"props"`
}

//statusPartition
type stPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

type statusGame struct {
	LobbyID int    `fesl:"lid"`
	Fit     int    `fesl:"fit"`
	GameID  string `fesl:"gid"`
}

type statusPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

// Status pnow.Status command
func (fm *FeslManager) Status(event network.EventClientCommand) {
	logrus.Println("==Status==")
	gameID := matchmaking.FindGIDs()

	ans := ansStatus{
		Txn:          pnowStatus,
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
