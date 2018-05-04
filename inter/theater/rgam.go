package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansRGAM struct {
	GameID  string `fesl:"GID"`
	LobbyID string `fesl:"LID"`
}

// RGAM - TEST
func (tm *Theater) RGAM(event network.EvProcess) {
	logrus.Println("=====RGAM REQUEST======")

	event.Client.Answer(&codec.Packet{
		Message: thtrRGAM,
		Content: ansRGAM{
			event.Process.Msg["GID"],
			event.Process.Msg["LID"],
		},
	})
}
