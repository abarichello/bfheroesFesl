package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)


// Lobbies Data
type ansHTSN struct {
	gameID   	string 	`fesl:"GID"`
	hostID 	 	string 	`fesl:"HPID"`
	UGID    	string 	`fesl:"UGID"`
	Secret   	string 	`fesl:"SECRET"`
	LobbyID     string  `fesl:"LID"`
}

// HTSN - TEst2
func (tm *Theater) HTSN(event network.EvProcess) {
		logrus.Println("HTSN HTSN===")
		// Client
		event.Client.Answer(&codec.Packet{
		Message: thtrHTSN,
		Content: ansHTSN{
			LobbyID: 	event.Process.Msg["LID"],
			Secret: 	"NOSECRET",
			gameID: 	event.Process.Msg["GID"],
			UGID: 		"NOGUID",
			hostID:     event.Process.Msg["HPID"],
		},
	})
}

