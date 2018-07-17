package fesl

import (
	"fmt"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ansNuLookupUserInfo struct {
	TXN      string     `fesl:"TXN"`
	UserInfo []userInfo `fesl:"userInfo"`
}

type userInfo struct {
	Namespace    string `fesl:"namespace"`
	XUID         int    `fesl:"xuid"`
	MasterUserID string `fesl:"masterUserId"`
	UserID       string `fesl:"userId"`
	UserName     string `fesl:"userName"`
	ClientID     string `fesl:"cid"`
}

func (fm *Fesl) NuLookupUserInfo(event network.EvProcess) {

	if event.Client.HashState.Get("clientType") == "server" && event.Process.Msg["userInfo.0.userName"] == "MargeSimpson" {
		fm.NuLookupUserInfoServer(event)
		return
	}

	answer := ansNuLookupUserInfo{
		TXN:      "NuLookupUserInfo",
		UserInfo: []userInfo{}}

	keys, _ := strconv.Atoi(event.Process.Msg["userInfo.[]"])
	for i := 0; i < keys; i++ {

		heroNamePkt := event.Process.Msg[fmt.Sprintf("userInfo.%d.userName", i)]

		var id, userID, heroName, online string
		err := fm.db.stmtGetHeroByName.QueryRow(heroNamePkt).Scan(&id, &userID, //br
			&heroName, &online) //auth
		if err != nil {
			return
		}
		logrus.Println("TEST TEST")
		answer.UserInfo = append(answer.UserInfo, userInfo{
			ClientID:     id,
			UserName:     heroName,
			UserID:       id,
			MasterUserID: id,
			Namespace:    "MAIN",
			XUID:         24,
		})
	}

	event.Client.Answer(&codec.Packet{
		Content: answer,
		Send:    event.Process.HEX,
		Message: acct,
	})

}

// NuLookupUserInfoServer - Server Login 1step
func (fm *Fesl) NuLookupUserInfoServer(event network.EvProcess) {
	var err error

	var id, userID, servername, secretKey, username string
	err = fm.db.stmtGetServerByID.QueryRow(event.Client.HashState.Get("sID")).Scan(&id, &userID, //br
		&servername, &secretKey, &username)

	if err != nil {
		logrus.Errorln(err)
		return
	}

	HEX := event.Process.HEX
	event.Client.Answer(&codec.Packet{
		Content: ansNuLookupUserInfo{
			TXN: "NuLookupUserInfo",
			UserInfo: []userInfo{
				{
					Namespace:    "MAIN",
					XUID:         24,
					MasterUserID: "1",
					UserID:       "1",
					ClientID:     "1",
					UserName:     servername,
				},
			},
		},
		Send:    HEX,
		Message: acct,
	})
}
