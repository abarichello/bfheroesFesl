package theater

import (
	"github.com/Synaxis/bfheroesFesl/backend/inter/network"
	"github.com/Synaxis/bfheroesFesl/backend/inter/network/codec"
)

type ansPING struct {
	TheaterID string `fesl:"TID"`
}

func (tm *Theater) PING(event *network.EventNewClient) {
	event.Client.WriteEncode(&codec.Packet{
		Type:    thtrPING,
		Payload: ansPING{"0"},
	})
}
