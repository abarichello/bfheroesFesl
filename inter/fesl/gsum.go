package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

type GetSessionId struct {
	TXN string `fesl:"TXN"`
}

func (fm *Fesl) gsumGetSessionID(event network.EvProcess) {
	event.Client.Answer(&codec.Packet{
		Content: GetSessionId{
			TXN: "GetSessionId"}, //pascal
			Message: "gsum",
	})
}

//fsys hello();