package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

type ansECHO struct {
	TID       	string `fesl:"TID"`
	TXN       	string `fesl:"TXN"`
	IP        	string `fesl:"IP"`
	Port      	int    `fesl:"PORT"`
	Error    	int    `fesl:"ERR"`
	TYPE      	int    `fesl:"TYPE"`
	UGID      	string `fesl:"UGID"`
	UID       	string `fesl:"UID"`


}

// ECHO - Possible ping from client to theater
func (tm *Theater) ECHO(event network.SocketUDPEvent) {
	Process := event.Data.(*network.ProcessFESL)
	ECHO := Process.Msg

	tm.socketUDP.Answer(&codec.Packet{
		Message: thtrECHO,
		Content: ansECHO{
			TID:       ECHO["TID"],
			TXN:       ECHO["TXN"],
			UGID:      ECHO["UGID"],
			TYPE:      1,
			UID:       ECHO["UID"],
			IP:        event.Addr.IP.String(),
			Port:      event.Addr.Port,
			Error: 0,						
		},
	}, event.Addr)
}