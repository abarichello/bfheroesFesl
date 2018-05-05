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

type Start struct {
	ID    	string              			`fesl:"id.id"`
	TXN  		string                   `fesl:"TXN"`
	Props   string 								   `fesl:"props.{}.[]"`
	idpart  string            			 `fesl:"id.partition"`

}

// Start handles pnow.Start
func (fm *Fesl) Start(event network.EvProcess) {
	logrus.Println("==START==")

	event.Client.Answer(&codec.Packet{
		Content: Start{
			ID: "1",
			TXN: "Start",
			idpart: partition,
		},
		Send:    event.Process.HEX,
		Message: pnow,
	})
	fm.Status(event)
}

