package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

type ansStart struct {
	TXN string      `fesl:"TXN"`
	ID  stPartition `fesl:"id"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientCommand) {
	event.Client.Answer(&codec.Pkt{
		Content: ansStart{
			TXN: pnowStart,
			ID: stPartition{1,
				event.Command.Msg[partition]},
		},
		Send: event.Command.HEX,
		Type: pnow,
	})
	fm.Status(event)
}
