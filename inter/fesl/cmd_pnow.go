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
	// pnowGetStatus = "GetStatus"
	pnowStart  = "Start"
	pnowStatus = "Status"
	// pnowUpdate    = "Update"
)

type ansStart struct {
	Taxon string          `fesl:"TXN"`
	ID    statusPartition `fesl:"id"`
}

//TODO: SYNC WITH DISCORD ROLE AND HWID
func (fm *FeslManager) Start(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("C Left")
		return
	}

	event.Client.WriteEncode(&codec.Packet{
		Payload: ansStart{
			Taxon: pnowStart,
			ID:    statusPartition{1, event.Command.Message["partition.partition"]},
		},
		Step: event.Command.PayloadID,
		Type: pnow,
	})

	fm.Status(event)
}

type ansStatus struct {
	Taxon        string                 `fesl:"TXN"`
	ID           statusPartition        `fesl:"id"`
	SessionState string                 `fesl:"sessionState"`
	Properties   map[string]interface{} `fesl:"props"`
}

type statusProperties struct {
	ResultType string       `fesl:"resultType"`
	Games      []statusGame `fesl:"games"`
}

type statusPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

type statusGame struct {
	LobbyID int    `fesl:"lit"`
	Fit     int    `fesl:"fit"`
	GameID  string `fesl:"gid"`
}

// Status - Basic fesl call to get overall service status (called before pnow?)
func (fm *FeslManager) Status(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("C Left")
		return
	}

	gameID := matchmaking.FindAvailableGIDs()
	
	ans := ansStatus{
		Taxon:  pnowStatus,
		ID: statusPartition{1, event.Command.Message["partition.partition"]},
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