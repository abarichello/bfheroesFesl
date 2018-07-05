package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

type ansGetTelemetryToken struct {
	Taxon          string `fesl:"TXN"`
	TelemetryToken string `fesl:"telemetryToken"`
	Enabled        bool   `fesl:"enabled"`
	Disabled       bool   `fesl:"disabled"`
}

// GetTelemetryToken
func (fm *Fesl) Telemetry(event network.EvProcess) {
	logrus.Println("Sent Telemetry")

	event.Client.Answer(&codec.Packet{
		Content: ansGetTelemetryToken{
			Taxon:          "GetTelemetryToken",
			TelemetryToken: `"teleToken"`,
			Enabled:        false,
			Disabled:       true,
		},

		Send:    event.Process.HEX,
		Message: acct,
	})
}
