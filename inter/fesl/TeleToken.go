package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

const (
	acctGetTelemetryToken = "GetTelemetryToken"
	// acctNuCreateEncryptedToken = "NuCreateEncryptedToken"
	// acctNuEntitleGame          = "NuEntitleGame"
	// acctNuEntitleUser          = "NuEntitleUser"
	// acctNuGetEntitlementCount  = "NuGetEntitlementCount"
	// acctNuGetEntitlements      = "NuGetEntitlements"
	// acctNuUpdateAccount        = "NuUpdateAccount"
	// acctTransactionException   = "TransactionException"
)

type ansGetTelemetryToken struct {
	Taxon          string `fesl:"TXN"`
	TelemetryToken string `fesl:"telemetryToken"`
	Filters        string `fesl:"filters"`
	Enabled        int 		`fesl:"enabled"`
	Disabled       int  	`fesl:"disabled"`
}

// GetTelemetryToken
func (fm *Fesl) Telemetry(event network.EvProcess) {
	logrus.Println("Sent Telemetry")

	event.Client.Answer(&codec.Packet{		
		Content: ansGetTelemetryToken{
			Taxon:          acctGetTelemetryToken,
			TelemetryToken: `"teleToken"`,
			Enabled:        0,
			Disabled: 			1,
		},

		Send: event.Process.HEX,
		Message: acct,
	})
}
