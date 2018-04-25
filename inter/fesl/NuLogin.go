package fesl

import (
	"github.com/satori/go.uuid"	
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)	
	
const (	
	acctNuLogin          = "NuLogin"
	acctNuLoginPersona   = "NuLoginPersona"
)

type userInfo struct {
	Namespace    string `fesl:"namespace"`
	XUID         string `fesl:"xuid"`
	MasterUserID string `fesl:"masterUserId"`
	UserID       string `fesl:"userId"`
	UserName     string `fesl:"userName"`
}

// Server Login Security -> Send close Packet
type NuLoginErr struct {
	TXN     string              `fesl:"TXN"`
	Message string              `fesl:"localizedMessage"`
	Errors  []LoginContainerErr `fesl:"errorContainer"`
	Code    int                 `fesl:"errorCode"`
}

type LoginContainerErr struct {
	Value      string `fesl:"value"`
	FieldError string `fesl:"fieldError"`
	FieldName  string `fesl:"fieldName"`
}

type ansNuLogin struct {
	TXN       string `fesl:"TXN"`
	ProfileID string `fesl:"profileId"`
	UserID    string `fesl:"userId"`
	NucleusID string `fesl:"nuid"`
	Encrypt   int   `fesl:"returnEncryptedInfo"`
	Lkey      string `fesl:"lkey"`
}

// NuLogin - First Login Command
func (fm *Fesl) NuLogin(event network.EvProcess) {

	if event.Client.HashState.Get("clientType") == "server" {
		// Server login
		fm.NuLoginServer(event)
		return
	}

	var id, username, email, birthday, language, country, gameToken string

	err := fm.db.stmtGetUserByGameToken.QueryRow(event.Process.Msg["encryptedInfo"]).Scan(&id, &username, //CONTINUE
		&email, &birthday, &language, &country, &gameToken) //todo add + checks 4 security

	if err != nil {
		event.Client.Answer(&codec.Packet{
			Content: NuLoginErr{
				TXN:     acctNuLogin,
				Message: `"Wrong Login/Spoof"`,
				Code:    120,
			},

			Send:    event.Process.HEX,
			Message: acct,
		})
		return
	}

	saveRedis := map[string]interface{}{
		"uID":       id,
		"username":  username,
		"sessionId": gameToken,
		"email":     email,
	}
	event.Client.HashState.SetM(saveRedis)

	// Setup a new key for our persona
	lkey := uuid.NewV4().String()
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", id)
	lkeyRedis.Set("userID", id)
	lkeyRedis.Set("name", username)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)
	event.Client.Answer(&codec.Packet{
		Content: ansNuLogin{
			TXN:       acctNuLogin,
			ProfileID: id,
			Encrypt:   1,
			UserID:    id,
			NucleusID: username,
			Lkey:      lkey,
		},
		Send:    event.Process.HEX,
		Message: acct,
	})
}

// NuLoginServer - login command for servers
func (fm *Fesl) NuLoginServer(event network.EvProcess) {

	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerBySecret.QueryRow(event.Process.Msg["password"]).Scan(&id,
		&userID, &servername, &secretKey, &username)

	if err != nil {
		event.Client.Answer(&codec.Packet{
			Content: NuLoginErr{
				TXN:     acctNuLogin,
				Message: `"Wrong Server "`,
				Code:    122,
			},
			Send:    event.Process.HEX,
			Message: acct,
		})
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



type reqNuLoginPersona struct {
	Txn  string `fesl:"TXN"`  // =NuLoginPersona
	Name string `fesl:"name"` // Value specified in +soldierName
}

type ansNuLoginPersona struct {
	TXN       string `fesl:"TXN"`
	ProfileID string `fesl:"profileId"`
	UserID    string `fesl:"userId"`
	Lkey      string `fesl:"lkey"`
}

// User log in with selected Hero
func (fm *Fesl) NuLoginPersona(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("C Left")
		return
	}

	if event.Client.HashState.Get("clientType") == "server" {
		logrus.Println("Server Login")
		fm.NuLoginPersonaServer(event)
		return
	}

	var id, userID, heroName, online string
	err := fm.db.stmtGetHeroeByName.QueryRow(event.Process.Msg["name"]).Scan(&id, &userID, &heroName, &online)
	if err != nil {
		logrus.Println("Wrong Login")
		return
	}

	// Setup a new key for our persona
	lkey := uuid.NewV4().String()
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", id)
	lkeyRedis.Set("userID", userID)
	lkeyRedis.Set("name", heroName)

	saveRedis := make(map[string]interface{})
	saveRedis["heroID"] = id
	event.Client.HashState.SetM(saveRedis)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)

	event.Client.Answer(&codec.Packet{
		Content: ansNuLogin{
			TXN:       acctNuLoginPersona,
			ProfileID: userID,
			UserID:    userID,
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
		logrus.Println("Server Login")
		return
	}

	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerByName.QueryRow(event.Process.Msg["name"]).Scan(&id, //continue
		&userID, &servername, &secretKey, &username)

	if event.Client.HashState.Get("clientType") != "server" {
		logrus.Println("Server Login Exploit")
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
		},
		Send:    event.Process.HEX,
		Message: acct,
	})
}