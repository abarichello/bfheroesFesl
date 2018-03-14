package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansEGRS struct {
	TID string `fesl:"TID"`
}

// EGRS - SERVER sent up, tell us if client is 'allowed' to join
func (tm *Theater) EGRS(event network.EventClientProcess) {
	if !event.Client.IsActive {
		return
	}
	logrus.Println("==EGRS==")

	if event.Process.Msg["ALLOWED"] == "1" {
		tm.db.stmtGameIncreaseJoining.Exec(event.Process.Msg["GID"])
	}

	event.Client.Answer(&codec.Packet{
		Message: thtrEGRS,
		Content: ansEGRS{event.Process.Msg["TID"]},
	})
}
