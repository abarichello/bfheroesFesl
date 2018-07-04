package mm

//Sync Theater with Fesl

import (
	 "github.com/Synaxis/bfheroesFesl/inter/network"
)

var Games = make(map[string]*network.Client)

// func FindGIDs() string {
// 	var gameID string	
// 	for id := range Games {
// 		gameID = id
// 		jsonStr := []string{id}
// 		logrus.WithFields(logrus.Fields{
// 			" ": jsonStr,
// 		}).Info("===Player Joined Game=== " + id) // TODO +uID joined Server
// 	}
// 	return gameID

// }
