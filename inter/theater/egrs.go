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

// EGRS - Enter Game Response
func (tm *Theater) EGRS(event network.EventClientProcess) {
	if !event.Client.IsActive {
		return
	}

	if event.Process.Msg["ALLOWED"] != "1" {
		return
	}

	logrus.Println("==EGRS==")
	tm.db.stmtGameIncreaseJoining.Exec(event.Process.Msg["GID"])

	event.Client.Answer(&codec.Packet{
		Message: thtrEGRS,
		Content: ansEGRS{
			event.Process.Msg["TID"],
			event.Process.Msg["PID"],
		},
	})
}
