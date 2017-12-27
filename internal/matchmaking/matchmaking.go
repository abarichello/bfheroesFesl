package matchmaking

import (
	"bitbucket.org/openheroes/backend/internal/network"
)

var Games = make(map[string]*network.Client)

func FindAvailableGIDs() string {
	var gameID string

	for k := range Games {
		gameID = k
	}

	return gameID
}
