package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

const (
	Nextchunk = "GetNextChunk"
)

type ansChunk struct {
	TXN   string `fesl:"TXN"`
	Chunk string `fesl:"chunk"`
	next  string `fesl:"GetNextChunk"`
}

func (fm *FeslManager) Chunk(event network.EventClientProcess) {
	hex := event.Process.HEX
	event.Client.Answer(&codec.Packet{
		Content: ansChunk{TXN: Nextchunk},
		Message: "Chunk",
		Send:    hex,
	})
}
