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
	TYPE   int    `fesl:"TYPE"`
	Send      string `fesl:"->D"`
	UGID       string `fesl:"UGID"`
	UID       string `fesl:"UID"`
	//SECRET    string `fesl:"SECRET"`


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
			Send:      ECHO["->D"],
			UGID:      ECHO["UGID"],
			TYPE:      1,
			UID:       ECHO["UID"],
			IP:        event.Addr.IP.String(),
			Port:      event.Addr.Port,
			ErrStatus: 0,						
		},
	}, event.Addr)
}