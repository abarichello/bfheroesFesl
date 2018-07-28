package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type reqUBRA struct {
	// TID=5
	TID int `fesl:"TID"`
	// LID=1
	LobbyID int `fesl:"LID"`
	// GID=3
	GameID int `fesl:"GID"`
	// START=0
	// START=1
	START int `fesl:"START"`
}

type ansUBRA struct {
	TID   string `fesl:"TID"`
	LID   string `fesl:"LID"`
	START int    `fesl:"START"`
}

// UBRA - "UpdateBracket" updates players connected (AP)
func (tM *Theater) UBRA(event network.EvProcess) {
	logrus.Println("=========UBRA===========")

	event.Client.Answer(&codec.Packet{
		Message: thtrUBRA,
		Content: ansUBRA{
			TID:   event.Process.Msg["TID"],
			LID:   event.Process.Msg["LID"],
			START: 1,
		}},
	)
	gdata := tM.level.NewObject("gdata", event.Process.Msg["GID"])	
	


	if event.Process.Msg["START"] == "1" {
		gdata.Set("AP", "0") // AP = Active-Player FOR THAT Player
		// If Player Entered -> Reset AP
	}

}
