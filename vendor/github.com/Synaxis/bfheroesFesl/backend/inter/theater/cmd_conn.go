package theater

import (
	"time"

	"github.com/Synaxis/bfheroesFesl/backend/inter/network"
	"github.com/Synaxis/bfheroesFesl/backend/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansClientConnected struct {
	TheaterID   string `fesl:"TID"`
	ConnectedAt int64  `fesl:"TIME"`
	ConnTTL     int    `fesl:"activityTimeoutSecs"`
	Protocol    string `fesl:"PROT"`
}

type reqCONN struct {
	Locale     string `fesl:"LOCALE"`     // "en_US",
	Platform   string `fesl:"PLAT"`       // "PC",
	Prod       string `fesl:"PROD"`       // "bfwest-pc",
	Protocol   int    `fesl:"PROT"`       // "2",
	SdkVersion string `fesl:"SDKVERSION"` // "5.0.0.0.0",
	Tid        int    `fesl:"TID"`        // "1",
	Version    string `fesl:"VERS"`       // "1.42.217478.0"
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
