package theater

import (
	"fmt"
	"net"

	"github.com/OSHeroes/bfheroesFesl/inter/mm"
	"github.com/OSHeroes/bfheroesFesl/inter/network"
	"github.com/OSHeroes/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

// ->N CGAM 0x40000000
type reqCGAM struct {
	// TID=3
	Tid int `fesl:"TID"`
	// LID=-1
	LobbyID int `fesl:"LID"`
	// RESERVE-HOST=0
	ReserveHost int `fesl:"RESERVE-HOST"`
	// NAME="[iad]A Battlefield Heroes Server(172.28.128.1:18567)"
	Name string `fesl:"NAME"`
	// PORT=18567
	Port int `fesl:"PORT"`
	// HTTYPE=A
	Httype string `fesl:"HTTYPE"`
	// TYPE=G
	Type string `fesl:"TYPE"`
	// QLEN=16
	Qlen int `fesl:"QLEN"`
	// DISABLE-AUTO-DEQUEUE=1
	DisableAutoDequeue int `fesl:"DISABLE-AUTO-DEQUEUE"`
	// HXFR=0
	Hxfr int `fesl:"HXFR"`
	// INT-PORT=18567
	IntPort int `fesl:"INT-PORT"`
	// INT-IP=192.168.1.102
	IntIP string `fesl:"INT-IP"`
	// MAX-PLAYERS=16
	MaxPlayers int `fesl:"MAX-PLAYERS"`
	// B-maxObservers=0
	BMaxObservers int `fesl:"B-maxObservers"`
	// B-numObservers=0
	BNumObservers int `fesl:"B-numObservers"`
	// UGID=GUID-Server
	Ugid string `fesl:"UGID"` /// Value passed in +guid
	// SECRET=Test-Server
	Secret string `fesl:"SECRET"` // Value passed in +secret
	// B-U-alwaysQueue=1
	BUAlwaysQueue int `fesl:"B-U-alwaysQueue"`
	// B-U-army_balance=Balanced
	BUArmyBalance string `fesl:"B-U-army_balance"`
	// B-U-army_distribution="0,0,0, 0,0,0,0, 0,0,0,0"
	BUArmyDistribution string `fesl:"B-U-army_distribution"`
	// B-U-avail_slots_national=yes
	BUAvailSlotsNational string `fesl:"B-U-avail_slots_national"`
	// B-U-avail_slots_royal=yes
	BUAvailSlotsRoyal string `fesl:"B-U-avail_slots_royal"`
	// B-U-avg_ally_rank=1000.0000
	BUAvgAllyRank string `fesl:"B-U-avg_ally_rank"`
	// B-U-avg_axis_rank=1000.0000
	BUAvgAxisRank string `fesl:"B-U-avg_axis_rank"`
	// B-U-community_name="Heroes SV"
	BUCommunityName string `fesl:"B-U-community_name"`
	// B-U-data_center=iad
	BUDataCenter string `fesl:"B-U-data_center"`
	// B-U-elo_rank=1000.0000
	BUEloRank string `fesl:"B-U-elo_rank"`
	// B-U-map=no_vehicles
	BUMap string `fesl:"B-U-map"`
	// B-U-percent_full=0
	BUPercentFull int `fesl:"B-U-percent_full"`
	// B-U-server_ip=172.28.128.1
	BUServerIP string `fesl:"B-U-server_ip"`
	// B-U-server_port=18567
	BUServerPort int `fesl:"B-U-server_port"`
	// B-U-server_state=empty
	BUServerState string `fesl:"B-U-server_state"`
	// B-version=1.46.222034.0
	BVersion string `fesl:"B-version"`
	// JOIN=O
	Join string `fesl:"JOIN"`
	// RT=
	Rt string `fesl:"RT"`
}
type ansCGAM struct {
	TID        string `fesl:"TID"`
	LobbyID    int `fesl:"LID"`
	MaxPlayers string `fesl:"MAX-PLAYERS"`
	EKEY       string `fesl:"EKEY"`
	UGID       string `fesl:"UGID"`
	Secret     string `fesl:"SECRET"`
	JOIN       string `fesl:"JOIN"`
	JoinMode   string `fesl:"JoindMode"`
	J          string `fesl:"J"`
	GameID     string `fesl:"GID"`
	isRanked   bool   `fesl:"B-U-UNRANKED"`
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
	gameServer.Set("QUEUE-LENGTH", "0")

	event.Client.HashState.Set("gdata:GID", gameID)

	_, err = tm.db.setServerStatsStatement(keys).Exec(args...)
	if err != nil {
		logrus.Error("Failed setting stats for game server "+gameID, err.Error())
		return
	}
	logrus.Println("===========CGAM=============")
	event.Client.Answer(&codec.Packet{
		Message: thtrCGAM,
		Content: ansCGAM{
			TID:        answer["TID"],
			LobbyID:    1, //should not be hardcoded
			UGID:       answer["UGID"],
			MaxPlayers: answer["MAX-PLAYERS"],
			EKEY:       "TEST1234",
			Secret:     "MargeSimpson",
			JOIN:       answer["JOIN"],
			isRanked:   false,
			J:          answer["JOIN"],
			JoinMode:   "1",
			GameID:     gameID,
		},
	})

	// Create game in database
	_, err = tm.db.stmtAddGame.Exec(gameID, addr.IP.String(), answer["PORT"], answer["B-version"], answer["JOIN"], answer["B-U-map"], 0, 0, answer["MAX-PLAYERS"], 0, 0, "")
	if err != nil {
		logrus.Println("Failed to add game: %v", err)
	}
	logrus.Println("Added GAMESERVER TO DB")
}
