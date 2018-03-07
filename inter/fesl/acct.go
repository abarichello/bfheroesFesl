package fesl

import (
	"strconv"

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
	acctNuLookupUserInfo = "NuLookupUserInfo"
	acctGrantEntitlement = "NuGrantEntitlement"
)

type ansNuLookupUserInfo struct {
	TXN      string     `fesl:"TXN"`
	UserInfo []userInfo `fesl:"userInfo"`
}

type userInfo struct {
	Namespace    string `fesl:"namespace"`
	XUID         string `fesl:"xuid"`
	MasterUserID string `fesl:"masterUserId"`
	UserID       string `fesl:"userId"`
	UserName     string `fesl:"userName"`
}

// NuLookupUserInfo - Gets basic information about a game user
func (fm *FeslManager) NuLookupUserInfo(event network.EventClientProcess) {

	if event.Client.HashState.Get("clientType") == "server" && event.Process.Msg["userInfo.0.userName"] == "Test-Server" {
		fm.NuLookupUserInfoServer(event)
		return
	}

	ans := ansNuLookupUserInfo{TXN: acctNuLookupUserInfo, UserInfo: []userInfo{}}

	logrus.Println("LookupUserInfo CLIENT" + event.Process.Msg["userInfo.0.userName"])

	keys, _ := strconv.Atoi(event.Process.Msg["userInfo.[]"])
	for i := 0; i < keys; i++ {
		heroNamePkt := event.Process.Msg["userInfo."+strconv.Itoa(i)+".userName"]

		var id, userID, heroName, online string
		err := fm.db.stmtGetHeroeByName.QueryRow(heroNamePkt).Scan(&id, &userID, &heroName, &online)
		if err != nil {
			return
		}

		ans.UserInfo = append(ans.UserInfo, userInfo{
			UserName:     heroName,
			UserID:       id,
			MasterUserID: id,
			Namespace:    "MAIN",
			XUID:         "24",
		})
	}

	event.Client.Answer(&codec.Pkt{
		Content: ans,
		Send:    event.Process.HEX,
		Type:    acct,
	})

}

// NuLookupUserInfoServer - Gets basic information about a game user
func (fm *FeslManager) NuLookupUserInfoServer(event network.EventClientProcess) {
	var err error

	var id, userID, servername, secretKey, username string
	err = fm.db.stmtGetServerByID.QueryRow(event.Client.HashState.Get("sID")).Scan(&id, &userID, &servername, &secretKey, &username)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	event.Client.Answer(&codec.Pkt{
		Content: ansNuLookupUserInfo{
			TXN: acctNuLookupUserInfo,
			UserInfo: []userInfo{
				{
					Namespace:    "MAIN",
					XUID:         "24",
					MasterUserID: "1",
					UserID:       "1",
					UserName:     servername,
				},
			},
		},
		Send: event.Process.HEX,
		Type: acct,
	})
}

type ansNuLoginPersona struct {
	TXN       string `fesl:"TXN"`
	ProfileID string `fesl:"profileId"`
	UserID    string `fesl:"userId"`
	Lkey      string `fesl:"lkey"`
}

// NuLoginPersona  // User logs in with selected Hero
func (fm *FeslManager) NuLoginPersona(event network.EventClientProcess) {
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
	lkey := BF2RandomUnsafe(24)
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", id)
	lkeyRedis.Set("userID", userID)
	lkeyRedis.Set("name", heroName)

	saveRedis := make(map[string]interface{})
	saveRedis["heroID"] = id
	event.Client.HashState.SetM(saveRedis)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)
	event.Client.Answer(&codec.Pkt{
		Content: ansNuLogin{
			TXN:       acctNuLoginPersona,
			ProfileID: userID,
			UserID:    userID,
			Lkey:      lkey,
		},
		Send: event.Process.HEX,
		Type: acct,
	})
}

//NuLoginPersonaServer Pre-Server Login (out of order ?)
func (fm *FeslManager) NuLoginPersonaServer(event network.EventClientProcess) {
	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerByName.QueryRow(event.Process.Msg["name"]).Scan(&id, //continue
		&userID, &servername, //continue
		&secretKey, &username)

	if err != nil {
		logrus.Println("Wrong Server Login")
		return
	}

	// Setup a new key for our persona
	lkey := BF2RandomUnsafe(24)
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", userID)
	lkeyRedis.Set("userID", userID)
	lkeyRedis.Set("name", servername)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)
	event.Client.Answer(&codec.Pkt{
		Content: ansNuLogin{
			TXN:       acctNuLoginPersona,
			ProfileID: id,
			UserID:    id,
			Lkey:      lkey,
		},
		Send: event.Process.HEX,
		Type: acct,
	})
}

type ansNuGetPersonas struct {
	TXN      string   `fesl:"TXN"`
	Personas []string `fesl:"personas"`
}

// NuGetPersonas . Display all Personas to the User
func (fm *FeslManager) NuGetPersonas(event network.EventClientProcess) {
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

	event.Client.Answer(&codec.Pkt{
		Send:    event.Process.HEX,
		Type:    acct,
		Content: ans,
	})
}

// test stuff
func (fm *FeslManager) NuGrantEntitlement(event network.EventClientProcess) {
	logrus.Println("GRANT ENTITLEMENT")

	event.Client.Answer(&codec.Pkt{
		Type:    acct,
		Content: "TXN",
		Send:    event.Process.HEX,
	})
}

// NuGetPersonasServer - Soldier data lookup call for servers
func (fm *FeslManager) NuGetPersonasServer(event network.EventClientProcess) {
	logrus.Println("SERVER CONNECT")

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
			logrus.Errorln(err)
			return
		}

		ans.Personas = append(ans.Personas, servername)
		event.Client.HashState.Set("ownerId."+strconv.Itoa(len(ans.Personas)), id)
	}

	event.Client.Answer(&codec.Pkt{
		Send:    event.Process.HEX,
		Type:    acct,
		Content: ans,
	})
}

// NuGetAccount - General account information retrieved, based on parameters sent
func (fm *FeslManager) NuGetAccount(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Client Left")
		return
	}

	fm.acctNuGetAccount(&event)
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
	ThidPartyOptIn bool   `fesl:"thidPartyOptin"`
}

func (fm *FeslManager) acctNuGetAccount(event *network.EventClientProcess) {
	event.Client.Answer(&codec.Pkt{
		Type: acct,
		Content: ansNuGetAccount{
			TXN:            acctNuGetAccount,
			Country:        "US",
			Language:       "en_US",
			DobDay:         1,
			DobMonth:       1,
			DobYear:        1992,
			GlobalOptIn:    false,
			ThidPartyOptIn: false,
			NucleusID:      event.Client.HashState.Get("email"),
			HeroName:       event.Client.HashState.Get("username"),
			UserID:         event.Client.HashState.Get("uID"),
		},
		Send: event.Process.HEX,
	})
}

type ansNuLogin struct {
	TXN       string `fesl:"TXN"`
	ProfileID string `fesl:"profileId"`
	UserID    string `fesl:"userId"`
	NucleusID string `fesl:"nuid"`
	Lkey      string `fesl:"lkey"`
}

type ansNuLoginErr struct {
	TXN     string                `fesl:"TXN"`
	Message string                `fesl:"localizedMessage"`
	Errors  []nuLoginContainerErr `fesl:"errorContainer"`
	Code    int                   `fesl:"errorCode"`
}

type nuLoginContainerErr struct {
	Value      string `fesl:"value"`
	FieldError string `fesl:"fieldError"`
	FieldName  string `fesl:"fieldName"`
}

// NuLogin - master login command
// TODO: Here we can implement a banlist/permission check if player is allowed to play/join
func (fm *FeslManager) NuLogin(event network.EventClientProcess) {

	if event.Client.HashState.Get("clientType") == "server" {
		// Server login
		fm.NuLoginServer(event)
		return
	}

	var id, username, email, birthday, language, country, gameToken string

	err := fm.db.stmtGetUserByGameToken.QueryRow(event.Process.Msg["encryptedInfo"]).Scan(&id, &username, //CONTINUE
		&email, &birthday, &language, &country, &gameToken)
	if err != nil {
		event.Client.Answer(&codec.Pkt{
			Content: ansNuLoginErr{
				TXN:     acctNuLogin,
				Message: `"Wrong Login/Spoof"`,
				Code:    120,
			},

			Send: event.Process.HEX,
			Type: event.Process.Query,
		})
		return
	}

	saveRedis := map[string]interface{}{
		"uID":       id,
		"username":  username,
		"sessionID": gameToken,
		"email":     email,
		"keyHash":   event.Process.Msg["encryptedInfo"],
	}
	event.Client.HashState.SetM(saveRedis)

	// Setup a new key for our persona
	lkey := BF2RandomUnsafe(24)
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", id)
	lkeyRedis.Set("userID", id)
	lkeyRedis.Set("name", username)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)
	event.Client.Answer(&codec.Pkt{
		Content: ansNuLogin{
			TXN:       acctNuLogin,
			ProfileID: id,
			UserID:    id,
			NucleusID: username,
			Lkey:      lkey,
		},
		Send: event.Process.HEX,
		Type: acct,
	})
}

// NuLoginServer - login command for servers
func (fm *FeslManager) NuLoginServer(event network.EventClientProcess) {
	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerBySecret.QueryRow(event.Process.Msg["password"]).Scan(&id,
		&userID, &servername, &secretKey, &username)

	if err != nil {
		event.Client.Answer(&codec.Pkt{
			Content: ansNuLoginErr{
				TXN:     acctNuLogin,
				Message: `"Wrong Server "`,
				Code:    122,
			},
			Send: event.Process.HEX,
			Type: acct,
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
	lkey := BF2RandomUnsafe(24)
	lkeyRedis := fm.level.NewObject("lkeys", lkey)
	lkeyRedis.Set("id", id)
	lkeyRedis.Set("userID", userID)
	lkeyRedis.Set("name", username)

	event.Client.HashState.Set("lkeys", event.Client.HashState.Get("lkeys")+";"+lkey)
	event.Client.Answer(&codec.Pkt{
		Content: ansNuLogin{
			TXN:       acctNuLogin,
			ProfileID: userID,
			UserID:    userID,
			NucleusID: username,
			Lkey:      lkey,
		},
		Send: event.Process.HEX,
		Type: acct,
	})
}
