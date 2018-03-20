package fesl

import (
	"strconv"

	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	rankGetStats    = "GetStats"
	rankUpdateStats = "UpdateStats"
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
func (fm *FeslManager) GetStats(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	answer := event.Process.Msg
	convert := strconv.Itoa
	owner := event.Process.Msg["owner"]
	userId := event.Client.HashState.Get("uID") //ultra typo

	if event.Client.HashState.Get("clientType") == "server" {

		var id, userID, heroName, online string
		err := fm.db.stmtGetHeroeByID.QueryRow(owner).Scan(&id, &userID, &heroName, &online)
		if err != nil {
			logrus.Println("ServerLOGIN")
			return
		}

		userId = userID
		logrus.Println("Server requesting stats")
	}

	ans := ansGetStats{
		TXN:       "GetStats",
		OwnerID:   owner,
		OwnerType: 1,
	}

	// Gen args list for statement -> heroID,userID,key1,key2,key3,..
	var args []interface{}
	statsKeys := make(map[string]string)
	args = append(args, owner)
	args = append(args, userId)
	keys, _ := strconv.Atoi(answer["keys.[]"])
	for i := 0; i < keys; i++ {
		args = append(args, answer["keys."+convert(i)+""])
		statsKeys[answer["keys."+convert(i)+""]] = convert(i)
	}

	rows, err := fm.db.getStatsStatement(keys).Query(args...)
	if err != nil {
		logrus.Errorln("Failed gettings stats for hero "+owner, err.Error())
	}

	for rows.Next() {
		var userID, heroID, statsKey, statsValue string
		err := rows.Scan(&userID, &heroID, &statsKey, &statsValue)
		if err != nil {
			logrus.Errorln("Issue with database:", err.Error())
		}

		ans.Stats = append(ans.Stats, statsPair{Key: statsKey, Value: statsValue, Text: statsValue})
		delete(statsKeys, statsKey)
	}

	// Send stats not found with value of ""
	for key := range statsKeys {
		ans.Stats = append(ans.Stats, statsPair{Key: key})
	}

	event.Client.Answer(&codec.Packet{
		Content: ans,
		Send:    event.Process.HEX,
		Message: "rank",
	})
}

type stat struct {
	text  string
	value float64
}

type ansUpdateStats struct {
	TXN   string      `fesl:"TXN"`
	Users []userStats `fesl:"u"`
}

type userStats struct {
	O     int          `fesl:"o"`
	Ot    int          `fesl:"ot"`
	Stats []updateStat `fesl:"s"`
}

type updateStat struct {
	Key   string `fesl:"k"`
	Pt    int    `fesl:"pt"`
	T     string `fesl:"t"`
	Ut    int    `fesl:"ut"`
	Value string `fesl:"v"`
}

// UpdateStats - updates stats about a soldier
func (fm *FeslManager) UpdateStats(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}
	answer := event.Process.Msg
	convert := strconv.Itoa
	ans := ansUpdateStats{TXN: rankUpdateStats, Users: []userStats{}}

	userId := event.Client.HashState.Get("uID")

	users, _ := strconv.Atoi(answer["u.[]"])

	if users == 0 {
		logrus.Warning("No u.[], defaulting to 1")
		users = 1
	}

	for i := 0; i < users; i++ {
		owner, ok := answer["u."+convert(i)+".o"]
		if event.Client.HashState.Get("clientType") == "server" {

			var id, userIDhero, heroName, online string
			err := fm.db.stmtGetHeroeByID.QueryRow(owner).Scan(&id, &userIDhero, &heroName, &online)
			if err != nil {
				logrus.Println("Persona not worthy!")
				return
			}
			if !ok { //check
				return
			}

			userId = userIDhero
			logrus.Println("Server updating stats")
		}

		if !ok { //check
			return
		}

		// Get current stats from DB
		// Make args list for the statement->heroID userID, key1, key2, key3,..
		stats := make(map[string]*stat)

		var argsGet []interface{}
		statsKeys := make(map[string]string)
		argsGet = append(argsGet, owner)
		argsGet = append(argsGet, userId)
		keys, _ := strconv.Atoi(answer["u."+convert(i)+".s.[]"])
		for j := 0; j < keys; j++ {
			argsGet = append(argsGet, answer["u."+convert(i)+".s."+convert(j)+".k"])
			statsKeys[answer["u."+convert(i)+".s."+convert(j)+".k"]] = convert(j)
		}

		rows, err := fm.db.getStatsStatement(keys).Query(argsGet...)
		if err != nil {
			logrus.Errorln("Failed gettings stats for hero "+owner, err.Error())
		}

		// Get all stats to be sent
		count := 0
		for rows.Next() {
			var userID, heroID, statsKey, statsValue string
			err := rows.Scan(&userID, &heroID, &statsKey, &statsValue)
			if err != nil {
				logrus.Errorln("Issue with database:", err.Error())
			}

			intValue, err := strconv.ParseFloat(statsValue, 64)
			if err != nil {
				intValue = 0
			}
			stats[statsKey] = &stat{
				text:  statsValue,
				value: intValue,
			}

			delete(statsKeys, statsKey)
			count++
		}

		if !event.Client.IsActive {
			logrus.Println("Cli Left")
			return
		}

		// Send stats not found with "" value
		for key := range statsKeys {
			stats[key] = &stat{
				text:  "",
				value: 0,
			}

			count++
		}
		// End getStats routine

		// Generate our argument list for the statement -> userId, owner, key1, value1, userId, owner, key2, value2, userId, owner, ...
		var args []interface{}
		keys, _ = strconv.Atoi(answer["u."+convert(i)+".s.[]"])
		for j := 0; j < keys; j++ {

			if answer["u."+convert(i)+".s."+convert(j)+".ut"] != "3" {
				logrus.Println("NewUpdate:", answer["u."+convert(i)+".s."+convert(j)+".k"], answer["u."+convert(i)+".s."+convert(j)+".t"], answer["u."+convert(i)+".s."+convert(j)+".ut"], answer["u."+convert(i)+".s."+convert(j)+".v"], answer["u."+convert(i)+".s."+convert(j)+".pt"])
			}

			key := answer["u."+convert(i)+".s."+convert(j)+".k"]
			value := answer["u."+convert(i)+".s."+convert(j)+".t"]

			if value == "" {
				logrus.Println("Updating stat", key+":", answer["u."+convert(i)+".s."+convert(j)+".v"], "+", stats[key].value)
				// We are dealing with a number
				value = answer["u."+convert(i)+".s."+convert(j)+".v"]

				// ut = 3 when we need to add up / if you level up = 0
				if answer["u."+convert(i)+".s."+convert(j)+".ut"] == "3" {
					intValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						// Couldn't transfer it to a number, skip updating this stat
						logrus.Errorln("Skipping stat "+key, err)
						event.Client.Answer(&codec.Packet{
							Send:    event.Process.HEX,
							Message: "rank",
							Content: ansUpdateStats{TXN: rankUpdateStats},
						})
						return
					}

					if !event.Client.IsActive {
						logrus.Println("Cli Left")
						return
					}

					if intValue <= 0 || event.Client.HashState.Get("clientType") == "server" || key == "c_ltp" || key == "c_sln" || key == "c_ltm" || key == "c_slm" || key == "c_wmid0" || key == "c_wmid1" || key == "c_tut" || key == "c_wmid2" {
						// limit keys for server only(TODO CHANGE THIS)
						newValue := stats[key].value + intValue

						if key == "c_wallet_hero" && newValue < 0 {
							logrus.Errorln("Negative STATS", key)
							event.Client.Answer(&codec.Packet{
								Send:    event.Process.HEX,
								Message: "rank",
								Content: ansUpdateStats{TXN: rankUpdateStats},
							})
							return
						}

						value = strconv.FormatFloat(newValue, 'f', 4, 64)
					} else {
						logrus.Errorln("Not allowed to process stat", key)
						event.Client.Answer(&codec.Packet{
							Send:    event.Process.HEX,
							Message: "rank",
							Content: ansUpdateStats{TXN: rankUpdateStats},
						})
						return
					}
				}
			}

			// We need to select 3 values for each insert/update,
			// owner, key and value
			logrus.Println("Updating stats:", userId, owner, key, value)
			args = append(args, userId)
			args = append(args, owner)
			args = append(args, key)
			args = append(args, value)
		}

		_, err = fm.db.setStatsStatement(keys).Exec(args...)
		if err != nil {
			logrus.Errorln("Failed setting stats for hero "+owner, err.Error())
		}
	}

	event.Client.Answer(&codec.Packet{
		Send:    event.Process.HEX,
		Message: "rank",
		Content: ans,
	})
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
func (fm *FeslManager) GetStatsForOwners(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	//refactor
	answer := event.Process.Msg
	convert := strconv.Itoa

	ans := ansGetStatsForOwners{
		TXN:   "rankGetStats", 
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
			err := fm.db.stmtGetHeroeByID.QueryRow(ownerID).Scan(&id, &userIDhero, &heroName, &online)
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
		logrus.Println("Cli Left")
		return
	}

	hex := event.Process.HEX
	event.Client.Answer(&codec.Packet{
		Send:    hex,
		Message: "rank",
		Content: ans,
	})
}
