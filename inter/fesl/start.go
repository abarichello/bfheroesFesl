package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/sirupsen/logrus"
)

//TODO
// 'GetStatus'
// 'Update'
// 'Cancel'

type Start struct {
	ID  stPartition `fesl:"id"`
	TXN string      `fesl:"TXN"`
}

// Start handles pnow.Start
func (fm *FeslManager) Start(event network.EventClientProcess) {
	logrus.Println("==START==")

	// var Button string

	// if event.Process.Msg["Start"] == "1" {
	// 	Button = "Start"
	// } else {
	// 	Button = "Cancel"
	// }

	// more info can be found at game files ./EABackend
	// 	PlayNowSettings.setPoolPlayerTimeout 30
	// PlayNowSettings.setPoolMaxPlayers 1
	// PlayNowSettings.enableEasyZone 0
	// PlayNowSettings.setEasyZoneLevel 5

	// rem -- Skill Level Search Parameter
	// PlayNowSettings.setEloWeight 250
	// PlayNowSettings.setEloFitThreshold 70
	// PlayNowSettings.setEloScale 210
	// PlayNowSettings.setEloScaleRetryDelta 210

	// rem -- Data Center Search Parameter
	// PlayNowSettings.setDataCenterWeight 500
	// PlayNowSettings.setDataCenterFitThreshold 50

	// rem -- Faction Balancing Search Parameter
	// PlayNowSettings.setFactionBalanceWeight 200
	// PlayNowSettings.setFactionBalanceFitThreshold 0

	// rem -- Server Capacity Search Parameter
	// PlayNowSettings.setPercentFullWeight 50
	// PlayNowSettings.setPercentFullFitThreshold 0
	// PlayNowSettings.setPercentFullScale 30
	// PlayNowSettings.setPercentFullTarget 80

	// rem - Epoch 0 => Find server with map preference and people in it (CLOSEST DATA CENTER ONLY)
	// rem - Epoch 1 => Find server with people in it (CLOSEST DATA CENTER ONLY)
	// PlayNowSettings.enableSearchEpoch 0 1

	// rem - Each tier will extend the EloScale by EloScaleRetryDelta
	// PlayNowSettings.enableTiersForSearchEpoch 0 1

	// rem == Fesl::PlayNowOptions::DebugThreshold
	// rem		0 = DEBUG_THRESHOLD_OFF
	// rem		1 = DEBUG_THRESHOLD_HIGH
	// rem		2 = DEBUG_THRESHOLD_MED
	// rem		3 = DEBUG_THRESHOLD_LOW
	// PlayNowSettings.setDebugLevel 1

	// PlayNowSettings.setSearchTimeout 60000
	// PlayNowSettings.setJoinTimeout 300000

	// rem PlayNowSettings.setMaxDataCenterPing 1000
	// rem PlayNowSettings.setEnableAutomaticFitTableGeneration 1
	// rem PlayNowSettings.setMaxAllowedLatency 250

	event.Client.Answer(&codec.Packet{
		Content: Start{
			TXN: "Start",
			ID: stPartition{1,
				event.Process.Msg[partition]},
		},
		Send:    event.Process.HEX,
		Message: "pnow",
	})
	fm.Status(event)
}
