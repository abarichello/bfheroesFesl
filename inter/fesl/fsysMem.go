package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

type reqMemCheck struct {
	// TXN stands for Taxon, sub-query name of the command.
	// Should be equal: MemCheck.
	TXN string `fesl:"TXN"`

	// FIXME: Result is usually an empty string
	Result string `fesl:"result"`
}

type ansMemCheck struct {
	// TXN stands for Taxon, sub-query name of the command.
	// Should be equal: MemCheck.
	TXN string `fesl:"TXN"`

	MemChecks []memCheck `fesl:"memcheck"`
	Salt      string     `fesl:"salt"`
	Type      int        `fesl:"type"`
}

type memCheck struct {
	Addr   string `fesl:"addr"`
	Length int    `fesl:"len"`
}

func (fm *Fesl) fsysMemCheck(event *network.EventNewClient) {	
	event.Client.Answer(&codec.Packet{
		Message: fsys,
		Content: ansMemCheck{
			TXN:      "MemCheck",
			Salt:     "",
		},
		Send: 0xC0000000,
	})
}