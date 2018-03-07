package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansEGRS struct {
	TheaterID string `fesl:"TID"`
}

// EGRS - SERVER sent up, tell us if client is 'allowed' to join
func (tm *Theater) EGRS(event network.EventClientProcess) {
	if !event.Client.IsActive {
		return
	}

	if event.Process.Msg["ALLOWED"] == "1" {
		_, err := tm.db.stmtGameIncreaseJoining.Exec(event.Process.Msg["GID"])
		if err != nil {
			logrus.Error("NOT Allowed ", err)
		}
	}

	event.Client.Answer(&codec.Pkt{
		Type:    thtrEGRS,
		Content: ansEGRS{event.Process.Msg["TID"]},
	})
}
