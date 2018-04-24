package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

// Lobbies Data
type ansLDAT struct {
	TID             string `fesl:"TID"`
	FavoriteGames   string `fesl:"FAVORITE-GAMES"`
	FavoritePlayers string `fesl:"FAVORITE-PLAYERS"`
	LobbyID         string `fesl:"LID"`
	Locale          string `fesl:"LOCALE"`
	MaxGames        string `fesl:"MAX-GAMES"`
	Name            string `fesl:"NAME"`
	NumGames        string `fesl:"NUM-GAMES"`
	Passing         string `fesl:"PASSING"`
}

func (tm *Theater) LobbyData(event network.EvProcess) {
	event.Client.Answer(&codec.Packet{
		Message: thtrLDAT,
		Content: ansLDAT{
			TID:             event.Process.Msg["TID"],
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

// Lobbies List
type ansLLST struct {
	TID        string `fesl:"TID"`
	NumLobbies int    `fesl:"NUM-LOBBIES"`
}

// LLST - CLIENT (???) unknown, potentially bookmarks
func (tm *Theater) LLST(event network.EvProcess) {
	event.Client.Answer(&codec.Packet{
		Message: thtrLLST,
		Content: ansLLST{event.Process.Msg["TID"], 1},
	})
}

// GLST - CLIENT called to get a list of game servers? Irrelevant for heroes.
func (tm *Theater) GLST(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}
	logrus.Println("GLST was called")
}
