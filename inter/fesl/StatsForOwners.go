package fesl

import (
	"strconv"

	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	rankGetStatsForOwners = "GetStatsForOwners"
)

type stat struct {
	text  string
	value float64
}

type ansGetStatsForOwners struct {
	TXN   string           `fesl:"TXN"`
	Stats []statsContainer `fesl:"stats"`
}

type statsContainer struct {
	Stats     []statsPair `fesl:"stats"`
	OwnerID   string      `fesl:"ownerId"`
	OwnerType int         `fesl:"ownerType"`
}

// GetStatsForOwners - Gives a bunch of info for the Hero selection screen?
func (fm *Fesl) GetStatsForOwners(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	//refactor
	answer := event.Process.Msg
	convert := strconv.Itoa

	ans := ansGetStatsForOwners{
		TXN:   "GetStats",
		Stats: []statsContainer{},
	}

	// Get the owner pids from redis
	numOfHeroes := event.Client.HashState.Get("numOfHeroes")
	userID := event.Client.HashState.Get("uID")
	numOfHeroesInt, err := strconv.Atoi(numOfHeroes)
	if err != nil {
		return
	}

	i := 1
	for i = 1; i <= numOfHeroesInt; i++ {
		ownerID := event.Client.HashState.Get("ownerId." + convert(i))
		if event.Client.HashState.Get("clientType") == "server" {

			var id, userIDhero, heroName, online string
			err := fm.db.stmtGetHeroByID.QueryRow(ownerID).Scan(&id, &userIDhero, &heroName, &online)
			if err != nil {
				logrus.Println("Weird getStats/Spoof")
				return
			}

			userID = userIDhero
			logrus.Println("===GetStats===")
		}

		stContainer := statsContainer{
			OwnerID:   ownerID,
			OwnerType: 1,
		}

		// Generate our argument list for the statement -> heroID, key1, key2, key3, ...
		var args []interface{}
		statsKeys := make(map[string]string)
		args = append(args, ownerID)
		args = append(args, userID)
		keys, _ := strconv.Atoi(answer["keys.[]"])
		for i := 0; i < keys; i++ {
			args = append(args, answer["keys."+convert(i)+""])
			statsKeys[answer["keys."+convert(i)+""]] = convert(i)
		}

		rows, err := fm.db.getStatsStatement(keys).Query(args...)
		if err != nil {
			logrus.Errorln("Failed gettings stats for hero "+ownerID, err.Error())
		}

		count := 0
		for rows.Next() {
			var userID, heroID, statsKey, statsValue string
			err := rows.Scan(&userID, &heroID, &statsKey, &statsValue)
			if err != nil {
				logrus.Errorln("Issue with database:", err.Error())
			}

			stContainer.Stats = append(stContainer.Stats, statsPair{
				Key:   statsKey,
				Value: statsValue,
				Text:  statsValue,
			})

			delete(statsKeys, statsKey)
			count++
		}

		for key := range statsKeys {
			stContainer.Stats = append(stContainer.Stats, statsPair{
				Key: key,
			})
		}

		ans.Stats = append(ans.Stats, stContainer)
	}

	if !event.Client.IsActive {
		logrus.Println("Client Left")
		return
	}
	if err != nil {
		logrus.Println("ERROR getStatsForOwners")
		return
	}

	hex := event.Process.HEX
	event.Client.Answer(&codec.Packet{
		Send:    hex,
		Message: "rank",
		Content: ans,
	})
}
