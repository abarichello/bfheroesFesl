package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansUBRA struct {
	TID string `fesl:"TID"`
	LID string `fesl:"LID"`
	START int `fesl:"START"`
}

// UBRA - "UpdateBracket" updates players connected
func (tM *Theater) UBRA(event network.EvProcess) {	
	logrus.Println("=========UBRA===========")
	//gdata := tM.level.NewObject("gdata", event.Process.Msg["GID"])

	event.Client.Answer(&codec.Packet{
		Message: thtrUBRA,
		Content: ansUBRA{
			TID: event.Process.Msg["TID"],
			LID: event.Process.Msg["LID"],
			START: 1,
		}},
	)

	// if event.Process.Msg["START"] == "0" {
	// 	gdata.Set("AP", "0")		
	// }

}
