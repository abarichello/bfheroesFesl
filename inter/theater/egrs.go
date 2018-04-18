package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansEGRS struct {
	TID string `fesl:"TID"`
	PID string `fesl:"PID"`
}

// EGRS - Enter Game Host Response
func (tm *Theater) EGRS(event network.EventClientProcess) {
	if !event.Client.IsActive {
		return
	}

	if event.Process.Msg["ALLOWED"] != "1" {
		// if ( !isAllowed )
 	   	// Fesl::Transaction::AddString(ft, "REASON", &reasonStr[4]);
		return
	}

	logrus.Println("======EGRS=====")
	tm.db.stmtGameIncreaseJoining.Exec(event.Process.Msg["GID"])

	event.Client.Answer(&codec.Packet{
		Message: thtrEGRS,
		Content: ansEGRS{
			event.Process.Msg["TID"],
			event.Process.Msg["PID"],
		},
	})
}

// GREM - Enter Game Host Response
func (tm *Theater) GREM(event network.EventClientProcess) {
	if !event.Client.IsActive {
		return
	}
	

	logrus.Println("======GREM=====")
	tm.db.stmtGameIncreaseJoining.Exec(event.Process.Msg["GID"])

	event.Client.Answer(&codec.Packet{
		Message: thtrEGRS,
		Content: ansEGRS{
			event.Process.Msg["LID"],
			event.Process.Msg["GID"],
		},
	})
}