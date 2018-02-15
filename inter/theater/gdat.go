package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

// GameClient Represents a game client connected to theater
type GameClient struct {
	ip   string
	port string
}

// GameServer Represents a game server and it's data
type GameServer struct {
	ip                 string
	port               string
	intIP              string
	intPort            string
	name               string
	level              string
	activePlayers      int
	maxPlayers         int
	queueLength        int
	joiningPlayers     int
	gameMode           string
	elo                float64
	numObservers       int
	maxObservers       int
	sguid              string
	hash               string
	password           string
	ugid               string
	sType              string
	join               string
	version            string
	dataCenter         string
	serverMap          string
	armyBalance        string
	armyDistribution   string
	availSlotsNational bool
	availSlotsRoyal    bool
	avgAllyRank        float64
	avgAxisRank        float64
	serverState        string
	communityName      string
}

type ansGDAT struct {
	Ap                  string `fesl:"AP"`
	ArmyDistribution    string `fesl:"B-U-army_distribution"`
	AvailableVipsNation string `fesl:"B-U-avail_vips_national"`
	AvailableVipsRoyal  string `fesl:"B-U-avail_vips_royal"`
	AvgAllyRank         string `fesl:"B-U-avg_ally_rank"`
	AvgAxisRank         string `fesl:"B-U-avg_axis_rank"`
	AvgLevel            string `fesl:"B-U-lvl_avg"`
	Easyzone            string `fesl:"B-U-easyzone"`
	EloRank             string `fesl:"B-U-elo_rank"`
	GameID              string `fesl:"GID"`
	IsRanked            string `fesl:"B-U-ranked"`
	Join                string `fesl:"JOIN"`
	LobbyID             string `fesl:"LID"`
	MapName             string `fesl:"B-U-map_name"`
	PunkBusterEnabled   string `fesl:"B-U-punkb"`
	ServerName          string `fesl:"NAME"`
	ServerType          string `fesl:"B-U-servertype"`
	StdDevLevel         string `fesl:"B-U-lvl_sdv"`
	TheaterID           string `fesl:"TID"`

	BMaxObservers        string `fesl:"B-maxObservers"`
	BNumObservers        string `fesl:"B-numObservers"`
	BUAlwaysQueue        string `fesl:"B-U-alwaysQueue"`
	BUArmyBalance        string `fesl:"B-U-army_balance"`
	BUAvailSlotsNational string `fesl:"B-U-avail_slots_national"`
	BUAvailSlotsRoyal    string `fesl:"B-U-avail_slots_royal"`
	BUCommunityName      string `fesl:"B-U-community_name"`
	BUDataCenter         string `fesl:"B-U-data_center"`
	BUMap                string `fesl:"B-U-map"`
	BUPercentFull        string `fesl:"B-U-percent_full"`
	BUServerIP           string `fesl:"B-U-server_ip"`
	BUServerPort         string `fesl:"B-U-server_port"`
	BUServerState        string `fesl:"B-U-server_state"`
	BVersion             string `fesl:"B-version"`
	DisableAutoDequeue   string `fesl:"DISABLE-AUTO-DEQUEUE"`
	Httype               string `fesl:"HTTYPE"`
	Hxfr                 string `fesl:"HXFR"`
	IntIp                string `fesl:"INT-IP"`
	IntPort              string `fesl:"INT-PORT"`
	IP                   string `fesl:"IP"`
	MaxPlayers           string `fesl:"MAX-PLAYERS"`
	Port                 string `fesl:"PORT"`
	Qlen                 string `fesl:"QLEN"`
	QueueLength          string `fesl:"QUEUE-LENGTH"`
	ReserveHost          string `fesl:"RESERVE-HOST"`
	Rt                   string `fesl:"RT"`
	Secret               string `fesl:"SECRET"`
	Type                 string `fesl:"TYPE"`
	Ugid                 string `fesl:"UGID"`
}

// GDAT - CLIENT called to get data about the server
func (tm *Theater) GDAT(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	gameID := event.Command.Message["GID"]
	gameServer := tm.level.NewObject("gdata", gameID)

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrGDAT,
		Payload: ansGDAT{
			TheaterID:            event.Command.Message["TID"],
			Ap:                   gameServer.Get("AP"),
			ArmyDistribution:     gameServer.Get("B-U-army_distribution"),
			AvailableVipsNation:  gameServer.Get("B-U-avail_vips_national"),
			AvailableVipsRoyal:   gameServer.Get("B-U-avail_vips_royal"),
			AvgAllyRank:          gameServer.Get("B-U-avg_ally_rank"),
			AvgAxisRank:          gameServer.Get("B-U-avg_axis_rank"),
			Easyzone:             gameServer.Get("B-U-easyzone"),
			EloRank:              gameServer.Get("B-U-elo_rank"),
			AvgLevel:             gameServer.Get("B-U-lvl_avg"),
			StdDevLevel:          gameServer.Get("B-U-lvl_sdv"),
			MapName:              gameServer.Get("B-U-map_name"),
			PunkBusterEnabled:    gameServer.Get("B-U-punkb"),
			IsRanked:             gameServer.Get("B-U-ranked"),
			ServerType:           gameServer.Get("B-U-servertype"),
			GameID:               gameServer.Get("GID"),
			Join:                 gameServer.Get("JOIN"),
			LobbyID:              gameServer.Get("LID"),
			ServerName:           gameServer.Get("NAME"),
			BMaxObservers:        gameServer.Get("B-maxObservers"),
			BNumObservers:        gameServer.Get("B-numObservers"),
			BUAlwaysQueue:        gameServer.Get("B-U-alwaysQueue"),
			BUArmyBalance:        gameServer.Get("B-U-army_balance"),
			BUAvailSlotsNational: gameServer.Get("B-U-avail_slots_national"),
			BUAvailSlotsRoyal:    gameServer.Get("B-U-avail_slots_royal"),
			BUCommunityName:      gameServer.Get("B-U-community_name"),
			BUDataCenter:         gameServer.Get("B-U-data_center"),
			BUMap:                gameServer.Get("B-U-map"),
			BUPercentFull:        gameServer.Get("B-U-percent_full"),
			BUServerIP:           gameServer.Get("B-U-server_ip"),
			BUServerPort:         gameServer.Get("B-U-server_port"),
			BUServerState:        gameServer.Get("B-U-server_state"),
			BVersion:             gameServer.Get("B-version"),
			DisableAutoDequeue:   gameServer.Get("DISABLE-AUTO-DEQUEUE"),
			Httype:               gameServer.Get("HTTYPE"),
			Hxfr:                 gameServer.Get("HXFR"),
			IntIp:                gameServer.Get("INT-IP"),
			IntPort:              gameServer.Get("INT-PORT"),
			IP:                   gameServer.Get("IP"),
			MaxPlayers:           gameServer.Get("MAX-PLAYERS"),
			Port:                 gameServer.Get("PORT"),
			Qlen:                 gameServer.Get("QLEN"),
			QueueLength:          gameServer.Get("QUEUE-LENGTH"),
			ReserveHost:          gameServer.Get("RESERVE-HOST"),
			Rt:                   gameServer.Get("RT"),
			Secret:               gameServer.Get("SECRET"),
			Type:                 gameServer.Get("TYPE"),
			Ugid:                 gameServer.Get("UGID"),
		},
	})
}
