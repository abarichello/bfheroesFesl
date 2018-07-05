package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type reqPLVT struct {
	// TID=16
	TID int `fesl:"TID"`

	// LID=1
	LobbyID int `fesl:"LID"`
	// GID=12
	GameID int `fesl:"GID"`
	// PID=3
	PlayerID int `fesl:"PID"`
}

type ansKICK struct {
	TID     string `fesl:"TID"`
	LobbyID string `fesl:"LID"`
	GID     string `fesl:"GID"`
}

type ansPLVT struct {
	TID      string `fesl:"TID"`
	PlayerID string `fesl:"PID"`
}

// PLVT - PlayerLeaveTeam
func (tM *Theater) PLVT(event network.EvProcess) {

	pid := event.Process.Msg["PID"]
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
			logrus.Println("Issue with database:", err.Error())
		}
		stats[statsKey] = statsValue
	}

	switch stats["c_team"] {
	case "1":
		_, err = tM.db.stmtGameDecreaseTeam1.Exec(event.Process.Msg["GID"])
		if err != nil {
			logrus.Println("PLVT ", err)
		}
	case "2":
		_, err = tM.db.stmtGameDecreaseTeam2.Exec(event.Process.Msg["GID"])
		if err != nil {
			logrus.Println("PLVT ", err)
		}
	default:
		logrus.Println("Invalid team " + stats["c_team"] + " for " + pid)
	}

	event.Client.Answer(&codec.Packet{ // need to check this
		Message: thtrPLVT,
		Content: ansPLVT{
			event.Process.Msg["PID"],
			event.Process.Msg["TID"],
		},
	})

	event.Client.Answer(&codec.Packet{ // need to check this
		Message: thtrKICK,
		Content: ansKICK{
			event.Process.Msg["PID"],
			event.Process.Msg["LID"],
			event.Process.Msg["TID"],
		},
	})

}
