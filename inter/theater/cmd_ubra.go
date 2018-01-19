package theater

import (
	"github.com/Synaxis/unstable/backend/inter/network"
	"github.com/Synaxis/unstable/backend/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansUBRA struct {
	TheaterID string `fesl:"TID"`
}

// UBRA - SERVER Called to  update server data
func (tM *Theater) UBRA(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrUBRA,
		Payload: ansUBRA{
			TheaterID: event.Command.Message["TID"],
		},
	})

	gdata := tM.level.NewObject("gdata", event.Command.Message["GID"])

	if event.Command.Message["START"] == "1" {
		gdata.Set("AP", "0")
	}

}
