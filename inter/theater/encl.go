package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansECNL struct {
	TheaterID string `fesl:"TID"`
	GameID    string `fesl:"GID"`
	LobbyID   string `fesl:"LID"`
}

// ECNL - CLIENT calls when they want to leave
func (tm *Theater) ECNL(event network.EventClientCommand) {
	logrus.Println("Left Q")

	event.Client.Answer(&codec.Packet{
		Step: 0x0,
		Type: thtrENCL,
		Payload: ansECNL{
			event.Command.Msg["TID"],
			event.Command.Msg["GID"],
			event.Command.Msg["LID"],
		},
	})
}
