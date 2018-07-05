package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	// Packet types

	FeslAccount     = "acct"
	FeslSystem      = "fsys"
	FeslGameSummary = "gsum"
	FeslPlayNow     = "pnow"
	FeslRanking     = "rank"
	FeslAssociation = "asso" // TODO: Not implemented
	FeslBlob        = "blob" // TODO: Not implemented
	FeslFeedback    = "fdbk" // TODO: Not implemented
	FeslFilter      = "fltr" // TODO: Not implemented
	FeslFindPlayer  = "fpla" // TODO: Not implemented
	FeslPresence    = "pres" // TODO: Not implemented
	FeslRecorder    = "recp" // TODO: Not implemented

	ThtrCreateGame           = "CGAM"
	ThtrConnect              = "CONN"
	ThtrEcho                 = "ECHO"
	ThtrEnterGame            = "EGAM"
	ThtrEnterGameEntitleGame = "EGEG"
	ThtrEnterGameRequest     = "EGRQ"
	ThtrEnterGameResponse    = "EGRS"
	ThtrEnterConnectionLost  = "ECNL"
	ThtrGamesData            = "GDAT"
	ThtrGamesList            = "GLST"
	ThtrKickPlayer           = "KICK"
	ThtrLobbyData            = "LDAT"
	ThtrLobbyList            = "LLST"
	ThtrPlayerEnter          = "PENT"
	ThtrPing                 = "PING"
	ThtrPlayerValidator      = "PLVT"
	ThtrUpdateBrokerRating   = "UBRA"
	ThtrUpdatePlayer         = "UPLA"
	ThtrUser                 = "USER"
)

type Packet struct {
	Message string
	Send    uint32
	Content interface{}
}

type RawPacket struct {
	// Query first 4 bytes
	// i.e. "fsys", "acct" "CONN", "UPLA"
	Query []byte

	// Broadcast, next 4 bytes
	// TOOD: Use enumerator
	Broadcast []byte

	// Length, next 4 bytes
	Length []byte

	// Payload
	Payload []byte
}

func ExtractPacket(buf *bytes.Buffer) (*RawPacket, error) {
	var err error

	// NOTE: .Data is nil right now
	pkt := &RawPacket{
		Query:     make([]byte, 4),
		Broadcast: make([]byte, 4),
		Length:    make([]byte, 4),
	}

	if _, err = buf.Read(pkt.Query); err != nil {
		return nil, err
	}
	if _, err = buf.Read(pkt.Broadcast); err != nil {
		return nil, err
	}
	if _, err = buf.Read(pkt.Length); err != nil {
		return nil, err
	}

	// 12 is number of actual read bytes 3x4=12
	portion := binary.BigEndian.Uint32(pkt.Length) - 12
	if portion <= 0 {
		return nil, fmt.Errorf("Undersized packet")
	}
	if int(portion) > buf.Len() {
		return nil, fmt.Errorf("Oversized packet")
	}

	pkt.Payload = make([]byte, int(portion))
	if _, err = buf.Read(pkt.Payload); err != nil {
		return nil, err
	}
	return pkt, nil
}
