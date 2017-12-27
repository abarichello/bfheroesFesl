package theater

import (
	"fmt"

	"bitbucket.org/openheroes/backend/internal/network"
	"bitbucket.org/openheroes/backend/internal/network/codec"
	"bitbucket.org/openheroes/backend/storage/level"

	"github.com/sirupsen/logrus"
)

type answerUSER struct {
	TheaterID string `fesl:"TID"`
	Name      string `fesl:"NAME"` // ServerName / ClientName
	ClientID  string `fesl:"CID"`  // ?
}

func (tm *Theater) NewState(ident string) *level.State {
	return tm.level.NewState(ident)
}

// USER - SHARED Called to get user data about client? No idea
func (tm *Theater) USER(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	lkeyRedis := tm.level.NewObject("lkeys", event.Command.Message["LKEY"])

	redisState := tm.NewState(fmt.Sprintf(
		"%s:%s",
		"mm",
		event.Command.Message["LKEY"],
	))
	event.Client.HashState = redisState

	redisState.Set("id", lkeyRedis.Get("id"))
	redisState.Set("userID", lkeyRedis.Get("userID"))
	redisState.Set("name", lkeyRedis.Get("name"))

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrUSER,
		Payload: answerUSER{
			TheaterID: event.Command.Message["TID"],
			Name:      lkeyRedis.Get("name"),
		},
	})
}
