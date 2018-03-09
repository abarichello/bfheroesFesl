package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

type ansPING struct {
	TID string `fesl:"TID"`
}

func (tm *Theater) PING(event *network.EventNewClient) {
	event.Client.Answer(&codec.Pkt{
		Type:    thtrPING,
		Content: ansPING{"0"},
	})
}
