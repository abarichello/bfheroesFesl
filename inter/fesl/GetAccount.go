package fesl

import (
	"github.com/OSHeroes/bfheroesFesl/inter/network"
	"github.com/OSHeroes/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	acct              = "acct"
	acctNuGetAccount  = "NuGetAccount"
	acctNuGetPersonas = "NuGetPersonas"
)

type ansNuGetAccount struct {
	TXN             string `fesl:"TXN"`
	NucleusID       string `fesl:"nuid"`
	UserID          string `fesl:"userId"`
	HeroName        string `fesl:"heroName"`
	DobDay          int    `fesl:"DOBDay"`
	DobMonth        int    `fesl:"DOBMonth"`
	DobYear         int    `fesl:"DOBYear"`
	Country         string `fesl:"country"`
	Language        string `fesl:"language"`
	GlobalOptIn     bool   `fesl:"globalOptin"`
	ThirdPartyOptIn bool   `fesl:"thirdPartyOptin"`
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
			TXN:             acctNuGetAccount,
			Country:         "US",
			Language:        "en_US",
			DobDay:          1,
			DobMonth:        1,
			DobYear:         1992,
			GlobalOptIn:     false,
			ThirdPartyOptIn: false,
			NucleusID:       event.Client.HashState.Get("email"),
			HeroName:        event.Client.HashState.Get("username"),
			UserID:          event.Client.HashState.Get("uID"),
		},
		Send: event.Process.HEX,
	})
}
