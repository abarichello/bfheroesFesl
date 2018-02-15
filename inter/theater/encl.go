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
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrENCL,
		Payload: ansECNL{
			event.Command.Message["TID"],
			event.Command.Message["GID"],
			event.Command.Message["LID"],
		},
	})
}
