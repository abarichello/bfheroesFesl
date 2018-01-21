package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansPENT struct {
	TheaterID string `fesl:"TID"`
	PlayerID  string `fesl:"PID"`
}

// PENT - SERVER sent up when a player joins (entitle player?)
func (tM *Theater) PENT(event network.EventClientCommand) {
	if !event.Client.IsActive {
		return
	}

	pid := event.Command.Message["PID"]

	// Get 4 stats for PID
	rows, err := tM.db.getStatsStatement(4).Query(pid, "c_kit", "c_team", "elo", "level")
	if err != nil {
		logrus.Errorln("Failed gettings stats for hero "+pid, err.Error())
	}

	stats := make(map[string]string)

	for rows.Next() {
		var userID, heroID, heroName, statsKey, statsValue string
		err := rows.Scan(&userID, &heroID, &heroName, &statsKey, &statsValue)
		if err != nil {
			logrus.Errorln("Issue with database:", err.Error())
		}
		stats[statsKey] = statsValue
	}

	switch stats["c_team"] {
	case "1":
		_, err = tM.db.stmtGameIncreaseTeam1.Exec(event.Command.Message["GID"])
		if err != nil {
			logrus.Error("PENT ", err)
		}
	case "2":
		_, err = tM.db.stmtGameIncreaseTeam2.Exec(event.Command.Message["GID"])
		if err != nil {
			logrus.Error("PENT ", err)
		}
	default:
		logrus.Errorln("Invalid team " + stats["c_team"] + " for " + pid)
	}

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrPENT,
		Payload: ansPENT{
			event.Command.Message["TID"],
			event.Command.Message["PID"],
		},
	})
}
