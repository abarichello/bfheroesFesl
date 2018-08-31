package theater

import (
	"strconv"

	"github.com/OSHeroes/bfheroesFesl/inter/network"

	"github.com/sirupsen/logrus"
)

type reqUPLA struct {
	// TID=12
	TID int `fesl:"TID"`
	// GID=3
	GameID int `fesl:"GID"`
	// LID=1
	LobbyID string `fesl:"LID"`
	// PID=6
	PlayerID int `fesl:"PID"`
	clientID int `fesl:"P-cid"`

	HostOwnerID int `fesl:"HMO"`
}

type reqUPLAKeys struct {
	PlayerElo   *string `fesl:"P-elo"`
	PlayerKills *string `fesl:"P-kills"`
	PlayerKit   *string `fesl:"P-kit"`
	PlayerLevel *string `fesl:"P-level"`
	// P-ping=24
	PlayerPing  *int    `fesl:"P-ping"`
	PlayerScore *string `fesl:"P-score"`
	PlayerTeam  *string `fesl:"P-team"`
	// P-time="1 min 10 sec "
	PlayerPlayedTime *string `fesl:"P-time"`
	PlayerClientID   *string `fesl:"P-cid"`
	PlayerDataCenter *string `fesl:"P-dc"`
	PlayerIP         *string `fesl:"P-ip"`
}

// Update Player
func (tM *Theater) UPLA(event network.EvProcess) {
	logrus.Println("==========UPLA==========")
	var args []interface{}

	keys := 0

	pid := event.Process.Msg["PID"]
	gid := event.Process.Msg["GID"]

	for index, value := range event.Process.Msg {
		if index == "TID" || index == "PID" || index == "GID" {
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

		args = append(args, gid)
		args = append(args, pid)
		args = append(args, index)
		args = append(args, value)
	}

	var err error
	_, err = tM.db.setServerPlayerStatsStatement(keys).Exec(args...)
	if err != nil {
		logrus.Println("Failed to update stats for player "+pid, err.Error())
	}

	gdata := tM.level.NewObject("gdata", event.Process.Msg["GID"])

	num, _ := strconv.Atoi(gdata.Get("AP"))

	num++

	gdata.Set("AP", strconv.Itoa(num))
}
