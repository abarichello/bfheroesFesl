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
	Type      int    `fesl:"TYPE"`
}

// ECHO - SHARED called like some heartbeat
func (tm *Theater) ECHO(event network.SocketUDPEvent) {
	command := event.Data.(*network.CommandFESL)

	tm.socketUDP.Answer(&codec.Packet{
		Type: thtrECHO,
		Payload: ansECHO{
			TXN:       command.Msg["TXN"],
			TID:       command.Msg["TID"],
			IP:        event.Addr.IP.String(),
			Port:      event.Addr.Port,
			ErrStatus: 0,
			Type:      1,
		},
	}, event.Addr)
}
