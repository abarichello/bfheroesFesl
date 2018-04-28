package theater

import (
	"net"
	"strconv"

	"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansEGAM struct {
	TID     string `fesl:"TID"`
	LobbyID string `fesl:"LID"`
	GameID  string `fesl:"GID"`
}

type ansEGRQ struct {
	TID          string `fesl:"TID"`
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
	Platform	 string `fesl:"PL"`
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
	TID      string `fesl:"TID"`
	Platform string `fesl:"PL"`
	Ticket   string `fesl:"TICKET"`
	PlayerID string `fesl:"PID"`
	IP       string `fesl:"I"`
	Port     string `fesl:"P"`
	Huid     string `fesl:"HUID"`
	Ekey     string `fesl:"EKEY"`
	IntIP    string `fesl:"INT-IP"`
	IntPort  string `fesl:"INT-PORT"`
	Secret   string `fesl:"SECRET"`
	Ugid     string `fesl:"UGID"`
	LobbyID  string `fesl:"LID"`
	GameID   string `fesl:"GID"`
}

// EGAM - EnterGameRequest
func (tm *Theater) EGAM(event network.EvProcess) {
	gameID := event.Process.Msg["GID"]
	externalIP := event.Client.IpAddr.(*net.TCPAddr).IP.String()
	lobbyID := event.Process.Msg["LID"]
	pid := event.Client.HashState.Get("id")  //playerID
	logrus.Println("====EGAM==")

	event.Client.Answer(&codec.Packet{
		Message: thtrEGAM,
		Content: ansEGAM{
			event.Process.Msg["TID"],
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

	if gameServer, ok := mm.Games[gameID]; ok {
		gsData := tm.level.NewObject("gdata", gameID)

		// Server
		gameServer.Answer(&codec.Packet{
			Message: thtrEGRQ,
			Content: ansEGRQ{
				TID:          event.Process.Msg["TID"],
				Name:         stats["heroName"],
				UserID:       stats["userID"],
				PlayerID:     pid,
				Ticket:       "2018751182",
				IP:           externalIP,
				Port:         strconv.Itoa(event.Client.IpAddr.(*net.TCPAddr).Port),
				IntIP:        event.Process.Msg["R-INT-IP"],
				IntPort:      event.Process.Msg["R-INT-PORT"],
				Ptype:        "P",
				RUser:        stats["heroName"],
				RUid:         stats["userID"],
				RUAccid:      stats["userID"],
				RUElo:        stats["elo"],
				RUTeam:       stats["c_team"],
				RUKit:        stats["c_kit"],
				RULvl:        stats["level"],
				RUDataCenter: "iad",
				Platform:	  event.Process.Msg["PC"],
				RUExternalIP: externalIP,
				RUInternalIP: event.Process.Msg["R-INT-IP"],
				RUCategory:   event.Process.Msg["R-U-category"],
				RIntIP:       event.Process.Msg["R-INT-IP"],
				RIntPort:     event.Process.Msg["R-INT-PORT"],
				Xuid:         "24",
				RXuid:        "24",
				LobbyID:      lobbyID,
				GameID:       gameID,
			},
		})
		logrus.Println("====EGRQ==")


		// Client
		event.Client.Answer(&codec.Packet{
			Message: thtrEGEG,
			Content: ansEGEG{
				TID:      event.Process.Msg["TID"],
				Ticket:   "2018751182",
				PlayerID: pid,
				IP:       gsData.Get("IP"),
				Port:     gsData.Get("PORT"),
				Huid:     "1",
				Ekey:     "O65zZ2D2A58mNrZw1hmuJw%3d%3d",
				IntIP:    gsData.Get("INT-IP"),
				IntPort:  gsData.Get("INT-PORT"),
				Secret:   "MargeSimpson",
				Platform: event.Process.Msg["PC"],
				Ugid:     gsData.Get("UGID"),
				LobbyID:  lobbyID,
				GameID:   gameID,
			},
		})
		logrus.Println("====EGEG==")

	}
}
