package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

const (
	gsum         = "gsum"
	GetSessionID = "GetSessionId"
)

type ansGetSessionID struct {
	TXN string `fesl:"TXN"`
}

func (fm *FeslManager) gsumGetSessionID(event network.EventClientProcess) {
	event.Client.Answer(&codec.Pkt{
		Content: ansGetSessionID{TXN: GetSessionID},
		Message: gsum,
	})
}
