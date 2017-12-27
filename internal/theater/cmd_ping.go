package theater

import (
	"bitbucket.org/openheroes/backend/internal/network"
	"bitbucket.org/openheroes/backend/internal/network/codec"
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
