package theater

import (
	"time"

	"github.com/Synaxis/unstable/backend/inter/network"
	"github.com/Synaxis/unstable/backend/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansClientConnected struct {
	TheaterID   string `fesl:"TID"`
	ConnectedAt int64  `fesl:"TIME"`
	ConnTTL     int    `fesl:"activityTimeoutSecs"`
	Protocol    string `fesl:"PROT"`
}

// CONN - SHARED (???) called on connection
func (tm *Theater) CONN(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrCONN,
		Payload: ansClientConnected{
			TheaterID:   event.Command.Message["TID"],
			ConnectedAt: time.Now().UTC().Unix(),
			ConnTTL:     int((60 * time.Minute).Seconds()),
			Protocol:    event.Command.Message["PROT"],
		},
	})
}
