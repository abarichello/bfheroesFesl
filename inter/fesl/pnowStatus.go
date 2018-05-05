package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	pnow 			= "pnow"
)

type Status struct {
	TXN  		string                   `fesl:"TXN"`
	ID    	string              `fesl:"id.id"`
	State 	string                   `fesl:"sessionState"`
	Props   string 								   `fesl:"props.{}.[]"`
	result    string 								 `fesl:"props.{resultType}"`
	idpart  string            			 `fesl:"id.partition"`
	PropsGames   map[string]interface{}                   `fesl:"props"`
}

type stGame struct {
	LobbyID int    `fesl:"lid"`
	Fit     int    `fesl:"fit"`
	GID     string `fesl:"gid"` //gameID
}


// Status comes after Start. tells info about desired server
func (fm *Fesl) Status(event network.EvProcess) {
	logrus.Println("=Status=")
	gameID := mm.FindGIDs()

	gamesArr := []stGame{
		{
			LobbyID: 1,
			GID:     gameID,
			Fit:     1001,
		},
	}

	//can't refactor more , please read go docs

	event.Client.Answer(&codec.Packet{
		Send:    0x80000000,
		Message: pnow,
		Content: Status{			
			TXN:   "Status",
			State: "COMPLETE",
			ID:    "1",
			idpart: partition,
			Props: "2",
			result: "JOIN",
			PropsGames: map[string]interface{}{
				"resultType": "JOIN",
				"games":      gamesArr,
			},
		}},
	)
}

// type Cancel struct {
// 	TXN   string                 `fesl:"TXN"`
// 	ID    stPartition            `fesl:"id"`
// 	State string                 `fesl:"sessionState"`
// 	Props map[string]interface{} `fesl:"props"`
// }

// // Cancel - cancel pnow
// func (fm *Fesl) Cancel(event network.EvProcess) {
// 	reply := event.Process.Msg

// 	event.Client.Answer(&codec.Packet{
// 		Content: Cancel{
// 			TXN:   "Cancel",
// 			State: "CANCELLED",
// 			ID:    stPartition{1, reply[partition]},
// 			Props: map[string]interface{}{
// 				"resultType": "CANCEL",
// 			},
// 		},
// 		Send:    0x80000000,
// 		Message: pnow,
// 	})
// }
