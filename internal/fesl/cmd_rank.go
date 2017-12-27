package fesl

import (
	"strconv"

	"bitbucket.org/openheroes/backend/internal/network"
	"bitbucket.org/openheroes/backend/internal/network/codec"

	"github.com/sirupsen/logrus"
)

const (
	rank = "rank"

	// rankGetDateRange            = "GetDateRange"
	// rankGetRankedStats          = "GetRankedStats"
	// rankGetRankedStatsForOwners = "GetRankedStatsForOwners"
	rankGetStats = "GetStats"
	// rankGetStatsForOwners       = "GetStatsForOwners"
	// rankGetTopN                 = "GetTopN"
	// rankGetTopNAndMe            = "GetTopNAndMe"
	// rankGetTopNAndStats         = "GetTopNAndStats"
	rankUpdateStats = "UpdateStats"
)

type ansGetStats struct {
	Taxon     string      `fesl:"TXN"`
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
func (fm *FeslManager) GetStats(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	owner := event.Command.Message["owner"]
	userId := event.Client.HashState.Get("uID")

	if event.Client.HashState.Get("clientType") == "server" {

		var id, userID, heroName, online string
		err := fm.db.stmtGetHeroeByID.QueryRow(owner).Scan(&id, &userID, &heroName, &online)
		if err != nil {
			logrus.Println("Persona not worthy!")
			return
		}

		userId = userID
		logrus.Println("Server requesting stats")
	}

	ans := ansGetStats{
		Taxon:     rankGetStats,
		OwnerID:   owner,
		OwnerType: 1,
	}

	// Generate our argument list for the statement -> heroID, userID, key1, key2, key3, ...
	var args []interface{}
	statsKeys := make(map[string]string)
	args = append(args, owner)
	args = append(args, userId)
	keys, _ := strconv.Atoi(event.Command.Message["keys.[]"])
	for i := 0; i < keys; i++ {
		args = append(args, event.Command.Message["keys."+strconv.Itoa(i)+""])
		statsKeys[event.Command.Message["keys."+strconv.Itoa(i)+""]] = strconv.Itoa(i)
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

	// Send stats not found with default value of ""
	for key := range statsKeys {
		ans.Stats = append(ans.Stats, statsPair{Key: key})
	}

	event.Client.WriteEncode(&codec.Packet{
		Payload: ans,
		Step:    event.Command.PayloadID,
		Type:    rank,
	})
}

type stat struct {
	text  string
	value float64
}

type ansUpdateStats struct {
	Taxon string      `fesl:"TXN"`
	Users []userStats `fesl:"u"`
}

// "u.0.o": "3",
// "u.0.ot": "1",
type userStats struct {
	O     int          `fesl:"o"`
	Ot    int          `fesl:"ot"`
	Stats []updateStat `fesl:"s"`
}

// "u.0.s.0.k": "c_ltp",
// "u.0.s.0.pt": "0",
// "u.0.s.0.t": "",
// "u.0.s.0.ut": "0",
// "u.0.s.0.v": "9025.0000",
type updateStat struct {
	Key   string `fesl:"k"`
	Pt    int    `fesl:"pt"`
	T     string `fesl:"t"`
	Ut    int    `fesl:"ut"`
	Value string `fesl:"v"`
}

// UpdateStats - updates stats about a soldier
func (fm *FeslManager) UpdateStats(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	ans := ansUpdateStats{Taxon: rankUpdateStats, Users: []userStats{}}

	userId := event.Client.HashState.Get("uID")

	users, _ := strconv.Atoi(event.Command.Message["u.[]"])

	if users == 0 {
		logrus.Warning("No u.[], defaulting to 1")
		users = 1
	}

	for i := 0; i < users; i++ {
		owner, ok := event.Command.Message["u."+strconv.Itoa(i)+".o"]
		if event.Client.HashState.Get("clientType") == "server" {

			var id, userIDhero, heroName, online string
			err := fm.db.stmtGetHeroeByID.QueryRow(owner).Scan(&id, &userIDhero, &heroName, &online)
			if err != nil {
				logrus.Println("Persona not worthy!")
				return
			}

			userId = userIDhero
			logrus.Println("Server updating stats")
		}

		if !ok {
			return
		}

		stats := make(map[string]*stat)

		// Get current stats from DB
		// Generate our argument list for the statement -> heroID, userID, key1, key2, key3, ...
		var argsGet []interface{}
		statsKeys := make(map[string]string)
		argsGet = append(argsGet, owner)
		argsGet = append(argsGet, userId)
		keys, _ := strconv.Atoi(event.Command.Message["u."+strconv.Itoa(i)+".s.[]"])
		for j := 0; j < keys; j++ {
			argsGet = append(argsGet, event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"])
			statsKeys[event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"]] = strconv.Itoa(j)
		}

		rows, err := fm.db.getStatsStatement(keys).Query(argsGet...)
		if err != nil {
			logrus.Errorln("Failed gettings stats for hero "+owner, err.Error())
		}

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

		// Send stats not found with default value of ""
		for key := range statsKeys {
			stats[key] = &stat{
				text:  "",
				value: 0,
			}

			count++
		}
		// end Get current stats from DB

		// Generate our argument list for the statement -> userId, owner, key1, value1, userId, owner, key2, value2, userId, owner, ...
		var args []interface{}
		keys, _ = strconv.Atoi(event.Command.Message["u."+strconv.Itoa(i)+".s.[]"])
		for j := 0; j < keys; j++ {

			if event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".ut"] != "3" {
				logrus.Println("Update new Type:", event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"], event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".t"], event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".ut"], event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".v"], event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".pt"])
			}

			key := event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".k"]
			value := event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".t"]

			if value == "" {
				logrus.Println("Updating stat", key+":", event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".v"], "+", stats[key].value)
				// We are dealing with a number
				value = event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".v"]

				// ut seems to be 3 when we need to add up (xp has ut 0 when you level'ed up, otherwise 3)
				if event.Command.Message["u."+strconv.Itoa(i)+".s."+strconv.Itoa(j)+".ut"] == "3" {
					intValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						// Couldn't transfer it to a number, skip updating this stat
						logrus.Errorln("Skipping stat "+key, err)
						event.Client.WriteEncode(&codec.Packet{
							Step:    event.Command.PayloadID,
							Type:    rank,
							Payload: ansUpdateStats{Taxon: rankUpdateStats},
						})
						return
					}

					if intValue <= 0 || event.Client.HashState.Get("clientType") == "server" || key == "c_ltp" || key == "c_sln" || key == "c_ltm" || key == "c_slm" || key == "c_wmid0" || key == "c_wmid1" || key == "c_tut" || key == "c_wmid2" {
						// Only allow increasing numbers (like HeroPoints) by the server for now
						newValue := stats[key].value + intValue

						if key == "c_wallet_hero" && newValue < 0 {
							logrus.Errorln("Not allowed to process stat. c_wallet_hero lower than 0", key)
							event.Client.WriteEncode(&codec.Packet{
								Step:    event.Command.PayloadID,
								Type:    rank,
								Payload: ansUpdateStats{Taxon: rankUpdateStats},
							})
							return
						}

						value = strconv.FormatFloat(newValue, 'f', 4, 64)
					} else {
						logrus.Errorln("Not allowed to process stat", key)
						event.Client.WriteEncode(&codec.Packet{
							Step:    event.Command.PayloadID,
							Type:    rank,
							Payload: ansUpdateStats{Taxon: rankUpdateStats},
						})
						return
					}
				}
			}

			// We need to append 3 values for each insert/update,
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

	event.Client.WriteEncode(&codec.Packet{
		Step:    event.Command.PayloadID,
		Type:    rank,
		Payload: ans,
	})
}

type ansGetStatsForOwners struct {
	Taxon string           `fesl:"TXN"`
	Stats []statsContainer `fesl:"stats"`
}

type statsContainer struct {
	Stats     []statsPair `fesl:"stats"`
	OwnerID   string      `fesl:"ownerId"`
	OwnerType int         `fesl:"ownerType"`
}

// GetStatsForOwners - Gives a bunch of info for the Hero selection screen?
func (fm *FeslManager) GetStatsForOwners(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	ans := ansGetStatsForOwners{
		Taxon: rankGetStats, // really? is it a typo? GetStatsForOwners?
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
		ownerID := event.Client.HashState.Get("ownerId." + strconv.Itoa(i))
		if event.Client.HashState.Get("clientType") == "server" {

			var id, userIDhero, heroName, online string
			err := fm.db.stmtGetHeroeByID.QueryRow(ownerID).Scan(&id, &userIDhero, &heroName, &online)
			if err != nil {
				logrus.Println("Persona not worthy!")
				return
			}

			userID = userIDhero
			logrus.Println("Server requesting stats")
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
		keys, _ := strconv.Atoi(event.Command.Message["keys.[]"])
		for i := 0; i < keys; i++ {
			args = append(args, event.Command.Message["keys."+strconv.Itoa(i)+""])
			statsKeys[event.Command.Message["keys."+strconv.Itoa(i)+""]] = strconv.Itoa(i)
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

	event.Client.WriteEncode(&codec.Packet{
		Step:    0xC0000007,
		Type:    rank,
		Payload: ans,
	})
}
