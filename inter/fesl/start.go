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

type Start struct {
	ID  stPartition `fesl:"id"`
	TXN string      `fesl:"TXN"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")

	   //  sub_8F25D8(&unk_121562C, (int)off_114BFD8, (int)"Start");
    // sub_8F25D8(&unk_1215614, (int)off_114BFD8, (int)"Cancel");
    // sub_8F25D8(&unk_1215638, (int)off_114BFD8, (int)"Update");
    // sub_8F25D8(&unk_1215620, (int)off_114BFD8, (int)"GetStatus");
    // sub_8F25D8(&unk_1215644, (int)off_114BFD8, (int)"Status");

	event.Client.Answer(&codec.Packet{
		Content: Start{
			TXN: "Start",
			ID: stPartition{1,
				event.Process.Msg[partition]},
		},
		Send:    event.Process.HEX,
		Message: "pnow",
	})
	fm.Status(event)
}
