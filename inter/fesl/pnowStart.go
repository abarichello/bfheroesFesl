package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)
//TODO
// 'GetStatus'
// 'Update'
// 'Cancel'
const (
	partition = "partition.partition"

)

type reqStart struct {
	// TXN=Start
	TXN string `fesl:"TXN"`
	// partition.partition=
	Partition string `fesl:"partition.partition"`
	// debugLevel=off
	debugLevel string `fesl:"debugLevel"`
	// version=1
	Version int `fesl:"version"`
	// players.[]=1
	//Players []reqStartPlayer
}

type Start struct {
	ID    				int                 `fesl:"id.id"`
	TXN  			    string              `fesl:"TXN"`
	Properties    		string 				`fesl:"props.{}.[]"`
  	Part  				string              `fesl:"id.partition"`

}

// Start handles pnow.Start
func (fm *Fesl) Start(event network.EvProcess) {
	logrus.Println("==START==")
	//var isSearching = true

	event.Client.Answer(&codec.Packet{
		Content: Start{
			TXN: "Start",
			ID: 	1,
			Part: event.Process.Msg[partition],
		},
		Send:    event.Process.HEX,
		Message: pnow,
	})
}

