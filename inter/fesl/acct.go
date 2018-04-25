package fesl

import (
	"strconv"
	uuid "github.com/satori/go.uuid"


	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	acct                 = "acct"
	acctNuGetAccount     = "NuGetAccount"
	acctNuGetPersonas    = "NuGetPersonas"
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
		// Server login
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
		logrus.Println("Client Left")
		return
	}

	if event.Client.HashState.Get("clientType") != "server" {
		// Server Exploit Login
		return
	}

	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerByName.QueryRow(event.Process.Msg["name"]).Scan(&id, //continue
		&userID, &servername, //continue
		&secretKey, &username)

	if event.Client.HashState.Get("clientType") != "server" {
		// Server Exploit Login
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

type ansNuGetPersonas struct {
	TXN      string   `fesl:"TXN"`
	Personas []string `fesl:"personas"`
}

// NuGetPersonas . Display all Personas/Heroes
func (fm *Fesl) NuGetPersonas(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Client Left")
		return
	}

	if event.Client.HashState.Get("clientType") == "server" {
		fm.NuGetPersonasServer(event)
		return
	}

	rows, err := fm.db.stmtGetHeroesByUserID.Query(event.Client.HashState.Get("uID"))
	if err != nil {
		return
	}

	ans := ansNuGetPersonas{TXN: acctNuGetPersonas, Personas: []string{}}

	for rows.Next() {
		var id, userID, heroName, online string
		err := rows.Scan(&id, &userID, &heroName, &online)
		if err != nil {
			logrus.Errorln(err)
			return
		}

		ans.Personas = append(ans.Personas, heroName)
		event.Client.HashState.Set("ownerId."+strconv.Itoa(len(ans.Personas)), id)
	}

	event.Client.HashState.Set("numOfHeroes", strconv.Itoa(len(ans.Personas)))

	event.Client.Answer(&codec.Packet{
		Send:    event.Process.HEX,
		Message: acct,
		Content: ans,
	})
}

// // test stuff
// func (fm *Fesl) NuGrantEntitlement(event network.EvProcess) {
// 	logrus.Println("GRANT ENTITLEMENT")

// 	event.Client.Answer(&codec.Packet{
// 		Message: "NuGrantEntitlement",
// 		Content: "TXN",
// 		Send:    event.Process.HEX,
// 	})
// }

// NuGetPersonasServer - Soldier data lookup call for servers
func (fm *Fesl) NuGetPersonasServer(event network.EvProcess) {
	logrus.Println("======SERVER CONNECTING=====")

	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerByName.QueryRow(event.Process.Msg["name"]).Scan(&id, //continue
		&userID, &servername, //continue
		&secretKey, &username)

	if event.Client.HashState.Get("clientType") != "server" {
		// Server Exploit Login
		logrus.Println("====Wrong Server Login====")

		return
	}

	// Server login
	rows, err := fm.db.stmtGetServerByID.Query(event.Client.HashState.Get("uID"))
	if err != nil {
		return
	}

	ans := ansNuGetPersonas{TXN: acctNuGetPersonas, Personas: []string{}}

	for rows.Next() {
		var id, userID, servername, secretKey, username string
		err := rows.Scan(&id, &userID, &servername, &secretKey, &username)
		if err != nil {
			event.Client.Answer(&codec.Packet{
				Content: NuLoginErr{
					TXN:     acctNuLogin,
					Message: `"Wrong Login/Spoof"`,
					Code:    120,
				},

				Send:    event.Process.HEX,
				Message: event.Process.Query,
			})
			return
		}

		ans.Personas = append(ans.Personas, servername)
		event.Client.HashState.Set("ownerId."+strconv.Itoa(len(ans.Personas)), id)
	}

	event.Client.Answer(&codec.Packet{
		Send:    event.Process.HEX,
		Message: "acct",
		Content: ans,
	})
}

type ansNuGetAccount struct {
	TXN            string `fesl:"TXN"`
	NucleusID      string `fesl:"nuid"`
	UserID         string `fesl:"userId"`
	HeroName       string `fesl:"heroName"`
	DobDay         int    `fesl:"DOBDay"`
	DobMonth       int    `fesl:"DOBMonth"`
	DobYear        int    `fesl:"DOBYear"`
	Country        string `fesl:"country"`
	Language       string `fesl:"language"`
	GlobalOptIn    bool   `fesl:"globalOptin"`
	ThirdPartyOptIn bool  `fesl:"thirdPartyOptin"`
}

// NuGetAccount - General account information retrieved, based on parameters sent
func (fm *Fesl) NuGetAccount(event network.EvProcess) {
	if !event.Client.IsActive {
		logrus.Println("Client Left")
		return
	}

	fm.acctNuGetAccount(&event)
}


func (fm *Fesl) acctNuGetAccount(event *network.EvProcess) {
	event.Client.Answer(&codec.Packet{
		Message: acct,
		Content: ansNuGetAccount{
			TXN:           		acctNuGetAccount,
			Country:        	"US",
			Language:       	"en_US",
			DobDay:         	1,
			DobMonth:       	1,
			DobYear:        	1992,
			GlobalOptIn:    	false,
			ThirdPartyOptIn:	false,
			NucleusID:      	event.Client.HashState.Get("email"),
			HeroName:       	event.Client.HashState.Get("username"),
			UserID:         	event.Client.HashState.Get("uID"),
		},
		Send: event.Process.HEX,
	})
}
