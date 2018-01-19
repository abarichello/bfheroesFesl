package theater

import (
	"net"
	"strconv"

	"github.com/Synaxis/unstable/backend/inter/matchmaking"
	"github.com/Synaxis/unstable/backend/inter/network"
	"github.com/Synaxis/unstable/backend/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansEGAM struct {
	TheaterID string `fesl:"TID"`
	LobbyID   string `fesl:"LID"`
	GameID    string `fesl:"GID"`
}

type ansEGRQ struct {
	TheaterID    string `fesl:"TID"`
	Name         string `fesl:"NAME"`
	UserID       string `fesl:"UID"`
	PlayerID     string `fesl:"PID"`
	Ticket       string `fesl:"TICKET"`
	IP           string `fesl:"IP"`
	Port         string `fesl:"PORT"`
	IntIP        string `fesl:"INT-IP"`
	IntPort      string `fesl:"INT-PORT"`
	Ptype        string `fesl:"PTYPE"`
	RUser        string `fesl:"R-USER"`
	RUid         string `fesl:"R-UID"`
	RUAccid      string `fesl:"R-U-accid"`
	RUElo        string `fesl:"R-U-elo"`
	RUTeam       string `fesl:"R-U-team"`
	RUKit        string `fesl:"R-U-kit"`
	RULvl        string `fesl:"R-U-lvl"`
	RUDataCenter string `fesl:"R-U-dataCenter"`
	RUExternalIP string `fesl:"R-U-externalIp"`
	RUInternalIP string `fesl:"R-U-internalIp"`
	RUCategory   string `fesl:"R-U-category"`
	RIntIP       string `fesl:"R-INT-IP"`
	RIntPort     string `fesl:"R-INT-PORT"`
	Xuid         string `fesl:"XUID"`
	RXuid        string `fesl:"R-XUID"`
	LobbyID      string `fesl:"LID"`
	GameID       string `fesl:"GID"`
}

type ansEGEG struct {
	TheaterID string `fesl:"TID"`
	Platform  string `fesl:"PL"`
	Ticket    string `fesl:"TICKET"`
	PlayerID  string `fesl:"PID"`
	IP        string `fesl:"I"`
	Port      string `fesl:"P"`
	Huid      string `fesl:"HUID"`
	Ekey      string `fesl:"EKEY"`
	IntIP     string `fesl:"INT-IP"`
	IntPort   string `fesl:"INT-PORT"`
	Secret    string `fesl:"SECRET"`
	Ugid      string `fesl:"UGID"`
	LobbyID   string `fesl:"LID"`
	GameID    string `fesl:"GID"`
}

// EGAM - CLIENT called when a client wants to join a gameserver
func (tm *Theater) EGAM(event network.EventClientCommand) {
	externalIP := event.Client.IpAddr.(*net.TCPAddr).IP.String()
	lobbyID := event.Command.Message["LID"]
	gameID := event.Command.Message["GID"]
	pid := event.Client.HashState.Get("id")

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrEGAM,
		Payload: ansEGAM{
			event.Command.Message["TID"],
			lobbyID,
			gameID,
		},
	})

	// Get 4 stats for PID
	rows, err := tm.db.getStatsStatement(4).Query(pid, "c_kit", "c_team", "elo", "level")
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

		stats["heroName"] = heroName
		stats["userID"] = userID
		stats[statsKey] = statsValue
	}

	if gameServer, ok := matchmaking.Games[gameID]; ok {
		gsData := tm.level.NewObject("gdata", gameID)

		// Server
		gameServer.WriteEncode(&codec.Packet{
			Type: thtrEGRQ,
			Payload: ansEGRQ{
				TheaterID:    "0",
				Name:         stats["heroName"],
				UserID:       stats["userID"],
				PlayerID:     pid,
				Ticket:       "2018751182",
				IP:           externalIP,
				Port:         strconv.Itoa(event.Client.IpAddr.(*net.TCPAddr).Port),
				IntIP:        event.Command.Message["R-INT-IP"],
				IntPort:      event.Command.Message["R-INT-PORT"],
				Ptype:        "P",
				RUser:        stats["heroName"],
				RUid:         stats["userID"],
				RUAccid:      stats["userID"],
				RUElo:        stats["elo"],
				RUTeam:       stats["c_team"],
				RUKit:        stats["c_kit"],
				RULvl:        stats["level"],
				RUDataCenter: "iad",
				RUExternalIP: externalIP,
				RUInternalIP: event.Command.Message["R-INT-IP"],
				RUCategory:   event.Command.Message["R-U-category"],
				RIntIP:       event.Command.Message["R-INT-IP"],
				RIntPort:     event.Command.Message["R-INT-PORT"],
				Xuid:         "24",
				RXuid:        "24",
				LobbyID:      lobbyID,
				GameID:       gameID,
			},
		})

		// Client
		event.Client.WriteEncode(&codec.Packet{
			Type: thtrEGEG,
			Payload: ansEGEG{
				TheaterID: event.Command.Message["TID"],
				Platform:  "pc",
				Ticket:    "2018751182",
				PlayerID:  pid,
				IP:        gsData.Get("IP"),
				Port:      gsData.Get("PORT"),
				Huid:      "1", // find via GID soon
				Ekey:      "O65zZ2D2A58mNrZw1hmuJw%3d%3d",
				IntIP:     gsData.Get("INT-IP"),
				IntPort:   gsData.Get("INT-PORT"),
				Secret:    "2587913",
				Ugid:      gsData.Get("UGID"),
				LobbyID:   lobbyID,
				GameID:    gameID,
			},
		})
	}
}
