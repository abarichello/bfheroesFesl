package fesl

import (
	"github.com/Synaxis/unstable/backend/inter/network"
	"github.com/Synaxis/unstable/backend/inter/network/codec"
)

const (
	gsum = "gsum"

	// gsumAddGameData            = "AddGameData"
	// gsumAddGameEvents          = "AddGameEvents"
	// gsumAddGameSummary         = "AddGameSummary"
	// gsumAddPlayerInfo          = "AddPlayerInfo"
	// gsumAddStats               = "AddStats"
	// gsumAddTeamInfo            = "AddTeamInfo"
	// gsumEndReport              = "EndReport"
	// gsumGameSummaryUpdateState = "GameSummaryUpdateState"
	// gsumGetGameData            = "GetGameData"
	// gsumGetGameEvents          = "GetGameEvents"
	// gsumGetGameHistory         = "GetGameHistory"
	// gsumGetGameSummary         = "GetGameSummary"
	// gsumGetPlayerInfo          = "GetPlayerInfo"
	gsumGetSessionID = "GetSessionId"
	// gsumGetTeamInfo            = "GetTeamInfo"
	// gsumStartReport            = "StartReport"
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
