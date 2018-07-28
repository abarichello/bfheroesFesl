package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/sirupsen/logrus"
	"strings"
)

type reqUGAM struct {
	// TID=14
	TID int `fesl:"TID"`

	// LID=1
	LobbyID int `fesl:"LID"`
	// GID=12
	GameID int `fesl:"GID"`
	// JOIN=O
	JoinMode string `fesl:"JOIN"`
	// MAX-PLAYERS=16
	MaxPlayers int `fesl:"MAX-PLAYERS"`
	// B-maxObservers=0
	MaxObservers int `fesl:"B-maxObservers"`
	// B-numObservers=0
	NumObservers int `fesl:"B-numObservers"`

	// reqUGAMKeys
}

type reqUGAMKeys struct {
	// B-U-army_balance=Axis
	// B-U-army_balance=Balanced
	// B-U-army_distribution="0,0,0,0,0,0,0,0,0,0,0"
	// B-U-army_distribution="1,0,0,1,1,0,0,0,0,0,0"
	// B-U-avail_vips_national=4
	// B-U-avail_vips_royal=4
	// B-U-avg_ally_rank=1000
	// B-U-avg_axis_rank=1000
	// B-U-easyzone=no
	// B-U-elo_rank=1000
	// B-U-lvl_avg=0.000000
	// B-U-lvl_sdv=0.000000
	// B-U-map_name=Village
	// B-U-map_name=seaside_skirmish
	// B-U-percent_full=0
	// B-U-percent_full=6
	// B-U-punkb=0
	// B-U-ranked=yes
	// B-U-server_state=empty
	// B-U-server_state=has_players
	// B-U-servertype=public
	// B-maxObservers=0
	// B-numObservers=0
	// NAME="[iad]A Battlefield Heroes Server(172.28.128.1:18567)"
}

func (tM *Theater) UGAM(event network.EvProcess) {
	
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
		logrus.Println("======UGAM  Error==== ", err)
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
