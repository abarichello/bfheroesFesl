package theater

import (
	"strings"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/sirupsen/logrus"
)

func (tM *Theater) UGAM(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}
	logrus.Println("==============UPDATE GAME==============")
	gameID := event.Process.Msg["GID"] // TODO gameID := mm.FindGids()

	gdata := tM.level.NewObject("gdata", gameID)

	logrus.Println("Updating GameServer " + gameID)

	var args []interface{}
	keys := 0
	for index, value := range event.Process.Msg {
		if index == "TID" {
			continue
		}

		keys++

		// Strip quotes
		value = strings.Trim(value, `"`)

		gdata.Set(index, value)
		args = append(args, gameID)
		args = append(args, index)
		args = append(args, value)
	}
	_, err := tM.db.stmtUpdateGame.Exec(gameID)
	if err != nil {
		logrus.Println("UGAM ", err)
	}

	_, err = tM.db.setServerStatsStatement(keys).Exec(args...)
	if err != nil {
		logrus.Println("Failed to update stats for game server "+gameID, err.Error())
		if err.Error() == "Error 1213: Deadlock found when trying to get lock; try restarting transaction" {
			_, err = tM.db.setServerStatsStatement(keys).Exec(args...)
			if err != nil {
				logrus.Println("Failed to update stats for game server "+gameID+" on the second try", err.Error())
			}
		}
	}
}
