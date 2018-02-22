package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansKICK struct {
	PlayerID string `fesl:"PID"`
	LobbyID  string `fesl:"LID"`
	GameID   string `fesl:"GID"`
}

type ansPLVT struct {
	TheaterID string `fesl:"TID"`
}

// PENT - SERVER sent up when a player joins (entitle player?)
func (tM *Theater) PLVT(event network.EventClientCommand) {
	if !event.Client.IsActive {
		return
	}

	pid := event.Command.Msg["PID"]

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
		_, err = tM.db.stmtGameDecreaseTeam1.Exec(event.Command.Msg["GID"])
		if err != nil {
			logrus.Error("PLVT ", err)
		}
	case "2":
		_, err = tM.db.stmtGameDecreaseTeam2.Exec(event.Command.Msg["GID"])
		if err != nil {
			logrus.Error("PLVT ", err)
		}
	default:
		logrus.Errorln("Invalid team " + stats["c_team"] + " for " + pid)
	}

	event.Client.Answer(&codec.Packet{
		Type: thtrKICK,
		Payload: ansKICK{
			event.Command.Msg["PID"],
			event.Command.Msg["LID"],
			event.Command.Msg["GID"],
		},
	})

	event.Client.Answer(&codec.Packet{
		Type:    thtrPLVT,
		Payload: ansPLVT{event.Command.Msg["TID"]},
	})
}
