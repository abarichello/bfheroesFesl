package theater

import (
	"bitbucket.org/openheroes/backend/internal/network"
	"bitbucket.org/openheroes/backend/internal/network/codec"

	"github.com/sirupsen/logrus"
)

type ansEGRS struct {
	TheaterID string `fesl:"TID"`
}

// EGRS - SERVER sent up, tell us if client is 'allowed' to join
func (tm *Theater) EGRS(event network.EventClientCommand) {
	if !event.Client.IsActive {
		return
	}

	if event.Command.Message["ALLOWED"] == "1" {
		_, err := tm.db.stmtGameIncreaseJoining.Exec(event.Command.Message["GID"])
		if err != nil {
			logrus.Error("EGRS ", err)
		}
	}

	event.Client.WriteEncode(&codec.Packet{
		Type:    thtrEGRS,
		Payload: ansEGRS{event.Command.Message["TID"]},
	})
}
