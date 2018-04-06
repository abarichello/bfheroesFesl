package mm

//match making

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/sirupsen/logrus"
)

var Games = make(map[string]*network.Client)

func FindGIDs() string {
	var gameID string

	for id := range Games {
		gameID = id
		jsonStr := []string{id}
		logrus.WithFields(logrus.Fields{
			" ": jsonStr,
		}).Info("===Player Joined Game=== " + id) // TODO +uID joined Server
	} //Synaxis joined Server 123
	return gameID
}
