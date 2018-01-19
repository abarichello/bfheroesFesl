package matchmaking

import (
	"github.com/Synaxis/unstable/backend/inter/network"
)

var Games = make(map[string]*network.Client)

func FindAvailableGIDs() string {
	var gameID string

	for k := range Games {
		gameID = k
	}

	return gameID
}
