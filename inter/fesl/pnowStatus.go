package fesl

import (
	"fmt"
	"encoding/binary"
	"net"
	"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"	
)

const (
	pnow = "pnow"
)

type Status struct {
	TXN  		     string     `fesl:"TXN"`
	ID    			 int        `fesl:"id.id"`
	State 			 string     `fesl:"sessionState"`
	Props   		 int 	    `fesl:"props.{}.[]"`
	result  		 string     `fesl:"props.{resultType}"`
	idpart  		 string     `fesl:"id.partition"`
	Properties   map[string]interface{} `fesl:"props"`
}

type stGame struct {
	LobbyID int    `fesl:"lid"`
	Fit     int    `fesl:"fit"`
	GID     string    `fesl:"gid"` //gameID to join
}


// Status comes after Start. tells info about desired server
func (fm *Fesl) Status(event network.EvProcess) {
	logrus.Println("=Status=")

	IP := binary.BigEndian.Uint32(event.Client.IpAddr.(*net.TCPAddr).IP.To4())
	gameID := mm.FindGIDs(event.Client.HashState.Get("heroID"), fmt.Sprint(IP))	

	// // continuos search
	// for GID := range gameID {
	// 	gamesArr := []stGame{
	// 	{
	// 		GID:     gameID,
	// 		Fit:     1001,
	// 		LobbyID: 1,
	// 	},
	// }

	event.Client.Answer(&codec.Packet{
		Send:    0x80000000,
		Message: event.Process.Query,
		Content: Status{			
			TXN:   "Status",
			State: "COMPLETE",
			ID:    1,
			idpart: event.Process.Msg[partition],
			Props: 2,
			result: "JOIN",
			Properties: map[string]interface{}{
				"gid": 	      gameID,
				"fit": 		  999,
				"lid":        1,
			},
		}},
	)
} 
