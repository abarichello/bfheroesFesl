package fesl

import (
	"strconv"

	"github.com/OSHeroes/bfheroesFesl/inter/network"
	"github.com/OSHeroes/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansNuGetPersonas struct {
	TXN      string   `fesl:"TXN"`
	Personas []string `fesl:"personas"`
}

// NuGetPersonas . Display all Personas/Heroes
func (fm *Fesl) NuGetPersonas(event network.EvProcess) {
	AFK := !event.Client.IsActive

	if AFK {
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

	ans := ansNuGetPersonas{
		TXN:      acctNuGetPersonas,
		Personas: []string{},
	}

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

// GetPersonasServer G_Server Login retrieves Info Based on +soldierName(should be more secure)
func (fm *Fesl) NuGetPersonasServer(event network.EvProcess) {
	logrus.Println("======SERVER CONNECT Prompt=====")
	//////Validates Login///////////
	AFK := !event.Client.IsActive
	if AFK {
		logrus.Println("Client Left")
		return
	}
	if event.Client.HashState.Get("clientType") != "server" {
		//Exploit Login
		logrus.Println("====Wrong Server Login====")
		return
	}

	var id, userID, servername, secretKey, username string

	err := fm.db.stmtGetServerByName.QueryRow(event.Process.Msg["name"]).Scan(&id, //continue
		&userID, &servername, &secretKey, &username)

	// Server login
	rows, err := fm.db.stmtGetServerByID.Query(event.Client.HashState.Get("uID"))
	if err != nil {
		return
	}

	//////Validates Login//////////////////////

	ans := ansNuGetPersonas{TXN: acctNuGetPersonas, Personas: []string{}}

	for rows.Next() {
		var id, userID, servername, secretKey, username string
		err := rows.Scan(&id, &userID, &servername, &secretKey, &username)
		if err != nil {
			logrus.Println("====Wrong Server Login====")
			return
		}

		ans.Personas = append(ans.Personas, servername)
		event.Client.HashState.Set("ownerId."+strconv.Itoa(len(ans.Personas)), id)
	}

	logrus.Println("====SERVER Login===")

	event.Client.Answer(&codec.Packet{
		Send:    event.Process.HEX,
		Message: acct,
		Content: ans,
	})
}
