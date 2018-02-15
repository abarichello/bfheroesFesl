package fesl

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

const (
	gsum = "gsum"

	// gsumAddGameData            = "AddGameData"
	// gsumAddGameEvents          = "AddGameEvents"
	// gsumAddGameSummary         = "AddGameSummary"
	// gsumAddPlayerInfo          = "AddPlayerInfo"
	// gsumAddStats               = "AddStats"
	// gsumAddTeamInfo            = "AddTeamInfo"
	// gsumGameSummaryUpdateState = "GameSummaryUpdateState"
	// gsumGetGameData            = "GetGameData"
	// gsumGetGameEvents          = "GetGameEvents"
	// gsumGetGameHistory         = "GetGameHistory"
	// gsumGetGameSummary         = "GetGameSummary"
	// gsumGetPlayerInfo          = "GetPlayerInfo"
	gsumGetSessionID = "GetSessionId"
	// gsumGetTeamInfo            = "GetTeamInfo"
)

type ansGetSessionID struct {
	Taxon string `fesl:"TXN"`
	// Games  []Game  `fesl:"games"`
	// Events []Event `fesl:"events"`
}

func (fm *FeslManager) gsumGetSessionID(event network.EventClientCommand) {
	event.Client.WriteEncode(&codec.Packet{
		Payload: ansGetSessionID{Taxon: gsumGetSessionID},
		Type:    gsum,
	})
}
