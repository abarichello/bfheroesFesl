package theater

import (
	"fmt"
	"net"

	"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansCGAM struct {
	TID        string `fesl:"TID"`
	LobbyID    string `fesl:"LID"`
	MaxPlayers string `fesl:"MAX-PLAYERS"`
	EKEY       string `fesl:"EKEY"`
	UGID       string `fesl:"UGID"`
	Secret     string `fesl:"SECRET"`
	JOIN       string `fesl:"JOIN"`
	JoinMode   string  `fesl:"JoindMode"`
	J          string `fesl:"J"`
	GameID     string `fesl:"GID"`
}

// CGAM - CreateGameParameters
func (tm *Theater) CGAM(event network.EvProcess) {

	answer := event.Process.Msg

	addr, ok := event.Client.IpAddr.(*net.TCPAddr)
	if !ok {
		logrus.Errorln("Failed turning IpAddr to net.TCPAddr")
		return
	}

	res, err := tm.db.stmtCreateServer.Exec(
		answer["NAME"],
		answer["B-U-community_name"],
		answer["INT-IP"],
		answer["INT-PORT"],
		answer["B-version"],
	)
	if err != nil {
		logrus.Error("Cannot create New server", err)
		return
	}

	id, _ := res.LastInsertId()
	gameID := fmt.Sprintf("%d", id)

	// Store gameID for access later
	mm.Games[gameID] = event.Client

	var args []interface{}

	// Setup a new key for our game
	gameServer := tm.level.NewObject("gdata", gameID)

	keys := 0

	// Stores what we know about this game in the redis db
	for index, value := range answer {
		if index == "TID" {
			continue
		}

		keys++

		// Strip quotes
		if len(value) > 0 && value[0] == '"' {
			value = value[1:]
		}
		if len(value) > 0 && value[len(value)-1] == '"' {
			value = value[:len(value)-1]
		}
		gameServer.Set(index, value)

		args = append(args, gameID)
		args = append(args, index)
		args = append(args, value)
	}

	gameServer.Set("LID", "1")
	gameServer.Set("GID", gameID)
	gameServer.Set("IP", addr.IP.String())
	gameServer.Set("AP", "0")
	gameServer.Set("QUEUE-LENGTH", "16")

	event.Client.HashState.Set("gdata:GID", gameID)

	_, err = tm.db.setServerStatsStatement(keys).Exec(args...)
	if err != nil {
		logrus.Error("Failed setting stats for game server "+gameID, err.Error())
		return
	}

	event.Client.Answer(&codec.Packet{
		Message: thtrCGAM,
		Content: ansCGAM{
			TID:        answer["TID"],
			LobbyID:    answer["1"],
			UGID:       answer["UGID"],
			MaxPlayers: answer["MAX-PLAYERS"],
			EKEY:       `O65zZ2D2A58mNrZw1hmuJw%3d%3d`,
			Secret:     `2587913`,
			JOIN:       answer["JOIN"],
			JoinMode: 	"1",
			J:          answer["J"],
			GameID:     gameID,
		},
	})
	logrus.Println("====CGAM====")


	// Create game in database
	_, err = tm.db.stmtAddGame.Exec(gameID, addr.IP.String(), answer["PORT"], answer["B-version"], answer["JOIN"], answer["B-U-map"], 0, 0, answer["MAX-PLAYERS"], 0, 0, "")
	if err != nil {
		logrus.Errorf("Failed to add game: %v", err)
	}
}
