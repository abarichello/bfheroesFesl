package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansEGRS struct {
	TID string `fesl:"TID"`
	PID string `fesl:"PID"`
	Allow string  `fesl:"ALLOWED"`
}

// EGRS - Enter Game Host Response
func (tm *Theater) EGRS(event network.EvProcess) {	

	if event.Process.Msg["ALLOWED"] == "1" {		
		return
	}

	logrus.Println("======EGRS=====")
	tm.db.stmtGameIncreaseJoining.Exec(event.Process.Msg["GID"])

	event.Client.Answer(&codec.Packet{
		Message: thtrEGRS,
		Content: ansEGRS{
			TID: event.Process.Msg["TID"],
			PID: event.Process.Msg["PID"],
			Allow: "1",
		},
	})
}

// Lobbies Data
type ansGREM struct {
	gameID   	string 	`fesl:"GID"`
	LID         string  `fesl:"LID"`
}

// GREM - Enter Game Host Response
func (tm *Theater) GREM(event network.EvProcess) {	

	logrus.Println("======GREM=====")
	event.Client.Answer(&codec.Packet{
		Message: thtrGREM,
		Content: ansGREM{
			event.Process.Msg["GID"],	
			event.Process.Msg["LID"],		
		},
	})
}