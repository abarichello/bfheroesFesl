package theater

import (
	"time"

	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansCONN struct {
	TID         string `fesl:"TID"`
	ConnectedAt int64  `fesl:"TIME"`
	ConnTTL     int    `fesl:"activityTimeoutSecs"`
	Protocol    string `fesl:"PROT"`
}

// CONN - Enters Theater
func (tm *Theater) CONN(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}
	
	
	logrus.Println("====CONN==")
	event.Client.Answer(&codec.Packet{
		Message: thtrCONN,
		Content: ansCONN{
			TID:         event.Process.Msg["TID"],
			ConnectedAt: time.Now().UTC().Unix(),
			ConnTTL:     int((60 * time.Minute).Seconds()),
			Protocol:    event.Process.Msg["PROT"],
		},
	})
}
