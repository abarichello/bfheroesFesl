package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansECNL struct {
	TID     string `fesl:"TID"`
	GameID  string `fesl:"GID"`
	LobbyID string `fesl:"LID"`
}

// ECNL - Called when Hero RQ/ Leave MM
func (tm *Theater) ECNL(event network.EventClientProcess) {
	logrus.Println("HeroRQ")

	event.Client.Answer(&codec.Packet{
		Message: thtrENCL,
		Content: ansECNL{
			event.Process.Msg["TID"],
			event.Process.Msg["GID"],
			event.Process.Msg["LID"],
		},
	})
}
