package matchmaking

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
)

var Games = make(map[string]*network.Client)

func FindGIDs() string {
	var gameID string

	for k := range Games {
		gameID = k
	}

	return gameID
}
