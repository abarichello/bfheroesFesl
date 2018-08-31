package theater

import (
	"time"

	"github.com/OSHeroes/bfheroesFesl/inter/network"
	"github.com/OSHeroes/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type reqCONN struct {
	// TID=1
	TID int `fesl:"TID"`

	// LOCALE=en_US
	Locale string `fesl:"LOCALE"`
	// PLAT=PC"
	Platform string `fesl:"PLAT"`
	// PROD=bfwest-pc
	Prod string `fesl:"PROD"`
	// PROT=2
	Protocol int `fesl:"PROT"`
	// SDKVERSION=5.0.0.0.0
	SdkVersion string `fesl:"SDKVERSION"`
	// VERS="1.42.217478.0 "
	Version string `fesl:"VERS"`
}

type ansCONN struct {
	TID         string `fesl:"TID"`
	ConnectedAt int64  `fesl:"TIME"`
	ConnTTL     int    `fesl:"activityTimeoutSecs"`
	Protocol    string `fesl:"PROT"`
}

// CONN - Enters Theater
func (tm *Theater) CONN(event network.EvProcess) {

	logrus.Println("======CONN=========")
	event.Client.Answer(&codec.Packet{
		Message: thtrCONN,
		Content: ansCONN{
			//sendPacket->SetVar("ATIME", "NuLoginPersona");
			TID:         event.Process.Msg["TID"],
			ConnectedAt: time.Now().UTC().Unix(),
			ConnTTL:     3600,
			Protocol:    event.Process.Msg["PROT"],
		},
	})
}
