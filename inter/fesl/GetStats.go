package fesl

import (
	"strconv"

	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	keysMult = "keys."
	keysArr  = "keys.[]"
	rank     = "rank"
)

type ansGetStats struct {
	TXN       string      `fesl:"TXN"`
	OwnerID   string      `fesl:"ownerId"`
	OwnerType int         `fesl:"ownerType"`
	Stats     []statsPair `fesl:"stats"`
}

type statsPair struct {
	Key   string `fesl:"key"`
	Text  string `fesl:"text"`
	Value string `fesl:"value"`
}

// GetStats - Get basic stats about a soldier/owner (account holder)
func (fm *Fesl) GetStats(event network.EvProcess) {

	answer := event.Process.Msg
	convert := strconv.Itoa


	owner := event.Process.Msg["owner"]
	userId := event.Client.HashState.Get("uID") //ultra typo

		if event.Client.HashState.Get("clientType") == "server" {
			logrus.Println("GetStats (server), replacing heroID with ownerID")
			var id, userID, heroName, online string
			err := fm.db.stmtGetHeroByID.QueryRow(owner).Scan(&id, &userID, &heroName, &online)
			if err != nil {
				logrus.Println("ServerLOGIN Error")
				return
			}

			userId = userID // should be userID = serverID (suID)
			//logrus.Println("Server requesting stats")

		}

	// Gen args list for statement -> heroID,userID,key1,key2,key3,..
	var args []interface{}
	statsKeys := make(map[string]string)
	args = append(args, owner)
	args = append(args, userId)
	keys, _ := strconv.Atoi(answer[keysArr])
	for i := 0; i < keys; i++ {
		args = append(args, answer[keysMult+convert(i)+""])
		statsKeys[answer[keysMult+convert(i)+""]] = convert(i)
	}

	ans := ansGetStats{
		TXN:       "GetStats",
		OwnerID:   owner,
		OwnerType: 1,
	}

	rows, err := fm.db.getStatsStatement(keys).Query(args...)

	// rows, err := fm.db.getStatsStatement(keys).Query(args...)
	if err != nil {
		logrus.Errorln("Failed gettings stats for hero "+owner, err.Error())
	}


	for rows.Next() {
		var userID, heroID, statsKey, statsValue string
		err := rows.Scan(&userID, &heroID, &statsKey, &statsValue)
		if err != nil {
			logrus.Errorln("Issue with GetStats:", err.Error())
		}

		ans.Stats = append(ans.Stats, statsPair{Key: statsKey, Value: statsValue, Text: statsValue})
		delete(statsKeys, statsKey)
	}

	// Send stats not found with default value of ""
	for key := range statsKeys {
		ans.Stats = append(ans.Stats, statsPair{Key: key})
	}
	
	event.Client.Answer(&codec.Packet{
		Content: ans,
		Send:    event.Process.HEX,
		Message: rank,
	})
}
