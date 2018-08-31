package theater

import (
	"github.com/OSHeroes/bfheroesFesl/inter/network"
	"github.com/OSHeroes/bfheroesFesl/inter/network/codec"
)

type ansECHO struct {
	TID   string `fesl:"TID"`
	IP    string `fesl:"IP"`
	Port  int    `fesl:"PORT"`
	Error int    `fesl:"ERR"`
	TYPE  int    `fesl:"TYPE"`
	UGID  string `fesl:"UGID"`
	TXN   string `fesl:"TXN"`
	UID   string `fesl:"UID"`
}

// ECHO - Broadcast The Server to NAT Interface
func (tm *Theater) ECHO(event network.SocketUDPEvent) {
	Process := event.Data.(*network.ProcessFESL)
	ECHO := Process.Msg

	tm.socketUDP.Answer(&codec.Packet{
		Message: thtrECHO,
		Content: ansECHO{
			TID:   ECHO["TID"],
			UGID:  ECHO["UGID"],
			TYPE:  1,
			Error: 0,
			UID:   ECHO["UID"],
			IP:    event.Addr.IP.String(),
			Port:  event.Addr.Port,
		},
	}, event.Addr)
}
