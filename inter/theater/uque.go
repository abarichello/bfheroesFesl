package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansUQUE struct {
	TID  string `fesl:"TID"`
	UQUE string `fesl:"UQUE"`
}

// UQUE - SERVER sent up, tell us if client is 'allowed' to join
func (tm *Theater) UQUE(event network.EvProcess) {
	if !event.Client.IsActive {
		return
	}
	logrus.Println("=====UQUE=======")

	event.Client.Answer(&codec.Packet{
		Message: thtrUQUE,
		Content: ansUQUE{
			event.Process.Msg["TID"],
			event.Process.Msg["UQUE"],
		},
	})
}
