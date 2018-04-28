package fesl

import (
	"github.com/satori/go.uuid"	
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

// NuLoginServer - NuLogin for gameServer.exe
func (fm *Fesl) NuLoginServer(event network.EvProcess) {

	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerBySecret.QueryRow(event.Process.Msg["password"]).Scan(&id,
		&userID, &servername, &secretKey, &username)

	if err != nil {
	logrus.Println("===NuLogin issue/wrong data!==")	
	return
	}

	saveRedis := make(map[string]interface{})
	saveRedis["uID"] = userID
	saveRedis["sID"] = id
	saveRedis["username"] = username
	saveRedis["apikey"] = event.Process.Msg["encryptedInfo"]
	saveRedis["keyHash"] = event.Process.Msg["password"]
	event.Client.HashState.SetM(saveRedis)


	// Setup a new key for new persona
	lkey := uuid.NewV4().String()
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", id)
	lkeyRedis.Set("userID", userID)
	lkeyRedis.Set("name", username)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)
	event.Client.Answer(&codec.Packet{
		Content: ansNuLogin{
			TXN:       acctNuLogin,
			ProfileID: userID,
			UserID:    userID,
			NucleusID: username,
			Lkey:      lkey,
		},
		Send:    event.Process.HEX,
		Message: acct,
	})
}

//NuLoginPersonaServer Pre-Server Login (out of order ?)
func (fm *Fesl) NuLoginPersonaServer(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("C Left")
		return
	}

	if event.Client.HashState.Get("clientType") != "server" {
		logrus.Println("======Possible Exploit=======")
		return
	}

	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerByName.QueryRow(event.Process.Msg["name"]).Scan(&id, //continue
		&userID, &servername, &secretKey, &username)

	if event.Client.HashState.Get("clientType") != "server" {
		logrus.Println("======Possible Exploit======")
		return
	}

	if err != nil {
		logrus.Println("Wrong Server Login")
		return
	}

	// Setup a key for Server
	lkey := uuid.NewV4().String()
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", userID)
	lkeyRedis.Set("userID", userID)
	lkeyRedis.Set("name", servername)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)
	event.Client.Answer(&codec.Packet{
		Content: ansNuLogin{
			TXN:       acctNuLoginPersona,
			ProfileID: id,
			UserID:    id,
			Lkey:      lkey,
			//nuid:      servername,
		},
		Send:    event.Process.HEX,
		Message: acct,
	})
}