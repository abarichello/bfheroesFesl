package theater

import (
	"fmt"

	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/Synaxis/bfheroesFesl/storage/level"

	"github.com/sirupsen/logrus"
)

type answerUSER struct {
	TID      string `fesl:"TID"`
	Name     string `fesl:"NAME"` // ServerName
	ClientID string `fesl:"CID"`  //clientID
}

func (tm *Theater) NewState(ident string) *level.State {
	return tm.level.NewState(ident)
}

// USER - Check User for Login - server not working atm
func (tm *Theater) USER(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	lkeyRedis := tm.level.NewObject("lkeys", event.Process.Msg["LKEY"])

	redisState := tm.NewState(fmt.Sprintf(
		"%s:%s",
		"mm",
		event.Process.Msg["LKEY"],
	))
	event.Client.HashState = redisState

	redisState.Set("id", lkeyRedis.Get("id"))
	redisState.Set("userID", lkeyRedis.Get("userID"))
	redisState.Set("name", lkeyRedis.Get("name"))

	event.Client.Answer(&codec.Packet{
		Message: thtrUSER,
		Content: answerUSER{
			ClientID: event.Process.Msg["CID"],
			TID:  event.Process.Msg["TID"],
			Name: lkeyRedis.Get("name"),
		},
	})
}
