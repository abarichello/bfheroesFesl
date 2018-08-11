package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

type GetSessionId struct {
	TXN string `fesl:"TXN"`
}

func (fm *Fesl) gsumGetSessionID(event network.EvProcess) {

	////Checks if we are a Server
	if event.Client.HashState.Get("clientType") == "server" {
		fm.NuLoginServer(event)
		return
	}

	event.Client.Answer(&codec.Packet{
		Content: GetSessionId{
			TXN: "GetSessionId", //yep , case sensitive
		},
		Message: "gsum",
	})
}
