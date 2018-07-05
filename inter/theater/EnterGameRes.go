package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type reqEGRQ struct {
	reqEGAM
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
	Platform     string `fesl:"PL"`
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

type reqEGRS struct {
	// TID=6
	TID int `fesl:"TID"`

	// LID=1
	LobbyID int `fesl:"LID"`
	// GID=12
	GameID int `fesl:"GID"`
	// ALLOWED=1
	// ALLOWED=0
	Allow int `fesl:"ALLOWED"`
	// PID=3
	PlayerID int `fesl:"PID"`

	// Reason is only sent when ALLOWED=0 and there is some kind of error
	// REASON=-602
	Reason string `fesl:"REASON,omitempty"`
}

type ansEGRS struct {
	TID   string `fesl:"TID"`
	LID   string `fesl:"LID"`
	PID   string `fesl:"PID"`
	Allow string `fesl:"ALLOWED"`
}

// EGRS - Enter Game Host Response
func (tm *Theater) EGRS(event network.EvProcess) {

	if event.Process.Msg["ALLOWED"] != "1" {
	}

	if event.Process.Msg["ALLOWED"] == "1" {
	}

	logrus.Println("======EGRS=====")
	tm.db.stmtGameIncreaseJoining.Exec(event.Process.Msg["GID"])

	event.Client.Answer(&codec.Packet{
		Message: thtrEGRS,
		Content: ansEGRS{
			TID:   event.Process.Msg["TID"],
			PID:   event.Process.Msg["PID"],
			LID:   event.Process.Msg["LID"],
			Allow: "1",
		},
	})
}

// Lobbies Data
type ansGREM struct {
	gameID string `fesl:"GID"`
	LID    string `fesl:"LID"`
}

// GREM - Enter Game Host Response
func (tm *Theater) GREM(event network.EvProcess) {

	logrus.Println("======GREM=====")
	event.Client.Answer(&codec.Packet{
		Message: thtrGREM,
		Content: ansGREM{
			event.Process.Msg["GID"],
			event.Process.Msg["LID"],
		},
	})
}
