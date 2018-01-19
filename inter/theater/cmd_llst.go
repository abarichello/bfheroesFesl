package theater

import (
	"github.com/Synaxis/unstable/backend/inter/network"
	"github.com/Synaxis/unstable/backend/inter/network/codec"
)

// Lobbies List
type ansLLST struct {
	TID        string `fesl:"TID"`
	NumLobbies int    `fesl:"NUM-LOBBIES"`
}

// Lobbies Data
type ansLDAT struct {
	TheaterID       string `fesl:"TID"`
	FavoriteGames   string `fesl:"FAVORITE-GAMES"`
	FavoritePlayers string `fesl:"FAVORITE-PLAYERS"`
	LobbyID         string `fesl:"LID"`
	Locale          string `fesl:"LOCALE"`
	MaxGames        string `fesl:"MAX-GAMES"`
	Name            string `fesl:"NAME"`
	NumGames        string `fesl:"NUM-GAMES"`
	Passing         string `fesl:"PASSING"`
}

// LLST - CLIENT (???) unknown, potentially bookmarks
func (tm *Theater) LLST(event network.EventClientCommand) {
	event.Client.WriteEncode(&codec.Packet{
		Type:    thtrLLST,
		Payload: ansLLST{event.Command.Message["TID"], 1},
	})

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrLDAT,
		Payload: ansLDAT{
			TheaterID:       "5",
			FavoriteGames:   "0",
			FavoritePlayers: "0",
			LobbyID:         "1",
			Locale:          "en_US",
			MaxGames:        "10000",
			Name:            "bfwestPC02",
			NumGames:        "1",
			Passing:         "0",
		},
	})
}
