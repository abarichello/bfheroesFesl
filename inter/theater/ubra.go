package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansUBRA struct {
	TID string `fesl:"TID"`
}

// UBRA - SERVER Called to  update server data
func (tM *Theater) UBRA(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	gdata := tM.level.NewObject("gdata", event.Process.Msg["GID"])

	if event.Process.Msg["Start"] == "1" {
		if event.Process.Msg["Status"] == "1" {
			gdata.Set("AP", "0")
		}
	}

	event.Client.Answer(&codec.Packet{
		Message: thtrUBRA,
		Content: ansUBRA{event.Process.Msg["TID"]},
	})
}
