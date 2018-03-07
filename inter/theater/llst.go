package theater

import (
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
)

// Lobbies List
type ansLLST struct {
	TID        string `fesl:"TID"`
	NumLobbies int    `fesl:"NUM-LOBBIES"`
}

// LLST - CLIENT (???) unknown, potentially bookmarks
func (tm *Theater) LLST(event network.EventClientProcess) {
	event.Client.Answer(&codec.Pkt{
		Type:    thtrLLST,
		Content: ansLLST{event.Process.Msg["TID"], 1},
	})
}
