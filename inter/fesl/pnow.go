package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	//"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/sirupsen/logrus"
)

const (
	partition =		 "partition.partition"
)
type stPartition struct {
	ID        int    `fesl:"id"`
	Partition string `fesl:"partition"`
}

type Start struct {
	ID  stPartition `fesl:"id"`
	TXN string      `fesl:"TXN"`
}

// Start handles pnow.Start
func (fm *Fesl) Start(event network.EvProcess) {
	logrus.Println("--START--")

	event.Client.Answer(&codec.Packet{
		Content: Start{TXN: "Start",
			ID: stPartition{1, "eagames/bfwest-dedicated"},
		},
		Send:    event.Process.HEX,
		Message: "pnow",
	})
}
type Status struct {
	TXN  				string    	`fesl:"TXN"`
	ID					int      	`fesl:"id.id"`
	State				string    	`fesl:"sessionState"`
	idpart  			string   	`fesl:"id.partition"`
	Props				int  		`fesl:"props.{}.[]"`
	Properties  		map[string]interface{} `fesl:"props"`
	result  			string    	`fesl:"props.{resultType}"`
}

type stGame struct {
	LobbyID 	int    	`fesl:"lid"`
	Fit     	int    	`fesl:"fit"`
	GID     	string 	`fesl:"gid"` //gameID to join
}


// Status comes after Start. tells info about desired server
func (fm *Fesl) Status(event network.EvProcess) {
	logrus.Println("--Status--")	

	var err error

	var gid string
	err = fm.db.stmtGetBookmark.QueryRow(event.Client.HashState.Get("uID")).Scan(&gid)
	if err != nil {
		return
	}

	gamesArray := []stGame{
		{
			GID:     gid,
			Fit:     1001,
			LobbyID: 1,
		},
	}

	// if event.Process.Msg["props.{games}.0.gid=0"] {
	// }

	event.Client.Answer(&codec.Packet{
		Send:    0x80000000,
		Message: event.Process.Query,
		Content: Status{
			TXN:   "Status",
			State: "COMPLETE",
			ID:    1,
			idpart: event.Process.Msg[partition],
			Props: 2,
			result: "JOIN",
			Properties: map[string]interface{}{
				"resultType": "JOIN",
				"games":      gamesArray,				
			}}},

	)
}//} //end for Loop
