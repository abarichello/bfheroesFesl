package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

type ansECHO struct {
	TID       string `fesl:"TID"`
	TXN       string `fesl:"TXN"`
	IP        string `fesl:"IP"`
	Port      int    `fesl:"PORT"`
	ErrStatus int    `fesl:"ERR"`
	Message   int    `fesl:"TYPE"`
}

//TODO check typo network.EventClientProcess
// ECHO - SHARED called like some heartbeat
func (tm *Theater) ECHO(event network.SocketUDPEvent) {
	Process := event.Data.(*network.ProcessFESL)
	ECHO := Process.Msg

	tm.socketUDP.Answer(&codec.Packet{
		Message: thtrECHO,
		Content: ansECHO{
			TID:       ECHO["TID"],
			TXN:       ECHO["TXN"],
			IP:        event.Addr.IP.String(),
			Port:      event.Addr.Port,
			ErrStatus: 0,
			Message:   1,
		},
	}, event.Addr)
}
