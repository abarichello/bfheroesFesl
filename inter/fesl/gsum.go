package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

const (
	gsum             = "gsum"
	gsumGetSessionID = "GetSessionId"
	// gsumAddGameData            = "AddGameData"
	// gsumAddStats               = "AddStats"
	// gsumGameSummaryUpdateState = "GameSummaryUpdateState"
	// gsumGetGameData            = "GetGameData"
	// gsumGetGameEvents          = "GetGameEvents"
	// gsumGetPlayerInfo          = "GetPlayerInfo"
)

type ansGetSessionID struct {
	Taxon string `fesl:"TXN"`
	// Games  []Game  `fesl:"games"`
	// Events []Event `fesl:"events"`
}

func (fm *FeslManager) gsumGetSessionID(event network.EventClientCommand) {
	event.Client.Answer(&codec.Packet{
		Payload: ansGetSessionID{
			Taxon: gsumGetSessionID},
		Type: gsum,
	})
}
