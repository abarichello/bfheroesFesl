package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansUBRA struct {
	TID string `fesl:"TID"`
	LID string `fesl:"LID"`
	// START string `fesl:"START"`


}

// UBRA - "UpdateBracket" updates players connected
func (tM *Theater) UBRA(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	gdata := tM.level.NewObject("gdata", event.Process.Msg["GID"])

	if event.Process.Msg["Start"] == "1" {
		gdata.Set("AP", "0")		
	}

	event.Client.Answer(&codec.Packet{
		Message: thtrUBRA,
		Content: ansUBRA{
			event.Process.Msg["TID"],
			event.Process.Msg["LID"],
		}})

}
