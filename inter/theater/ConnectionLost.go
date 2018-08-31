package theater

import (
	"github.com/OSHeroes/bfheroesFesl/inter/network"
	"github.com/OSHeroes/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type reqECNL struct {
	// TID=8
	TID int `fesl:"TID"`
	// GID=3
	GameID int `fesl:"GID"`
	// LID=1
	LobbyID int `fesl:"LID"`
}

type ansECNL struct {
	TID     string `fesl:"TID"`
	GameID  string `fesl:"GID"`
	LobbyID string `fesl:"LID"`
}

// ECNL - EnterConnectionLost
func (tm *Theater) ECNL(event network.EvProcess) {
	logrus.Println("=============ECNL==============")
	logrus.Println("========SENT Leave Announcement to Player======")
	logrus.Println("HeroRQ")

	event.Client.Answer(&codec.Packet{
		Message: thtrECNL,
		Content: ansECNL{
			event.Process.Msg["TID"],
			event.Process.Msg["GID"],
			event.Process.Msg["LID"],
		},
	})
}
