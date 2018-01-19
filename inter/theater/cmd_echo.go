package theater

import (
	"github.com/Synaxis/unstable/backend/inter/network"
	"github.com/Synaxis/unstable/backend/inter/network/codec"
)

type ansECHO struct {
	TheaterID string `fesl:"TID"`
	Taxon     string `fesl:"TXN"`
	IP        string `fesl:"IP"`
	Port      int    `fesl:"PORT"`
	ErrStatus int    `fesl:"ERR"`
	Type      int    `fesl:"TYPE"`
}

// ECHO - SHARED called like some heartbeat
func (tm *Theater) ECHO(event network.SocketUDPEvent) {
	command := event.Data.(*network.CommandFESL)

	tm.socketUDP.WriteEncode(&codec.Packet{
		Type: thtrECHO,
		Payload: ansECHO{
			Taxon:     command.Message["TXN"],
			TheaterID: command.Message["TID"],
			IP:        event.Addr.IP.String(),
			Port:      event.Addr.Port,
			ErrStatus: 0,
			Type:      1,
		},
	}, event.Addr)
}
