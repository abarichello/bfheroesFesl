package fesl

import (
	"fmt"
	"time"

	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

const (
	fsys = "fsys"

	fsysGetPingSites = "GetPingSites"
	// fsysGoodbye      = "Goodbye"
	fsysHello    = "Hello"
	fsysMemCheck = "MemCheck"
	fsysPing     = "Ping"
	// fsysSuicide      = "Suicide"
)

type ansMemCheck struct {
	Taxon     string     `fesl:"TXN"`
	MemChecks []memCheck `fesl:"memcheck"`
	Salt      string     `fesl:"salt"`
}

type memCheck struct {
	Length int    `fesl:"len"`
	Addr   string `fesl:"addr"`
}

func (fm *FeslManager) fsysMemCheck(event *network.EventNewClient) {
	event.Client.WriteEncode(&codec.Packet{
		Type: fsys,
		Payload: ansMemCheck{
			Taxon: fsysMemCheck,
			Salt:  "5",
		},
		Step: 0xC0000000,
	})
}

type ansHello struct {
	Taxon         string          `fesl:"TXN"`
	Domain        domainPartition `fesl:"domainPartition"`
	ConnTTL       int             `fesl:"activityTimeoutSecs"`
	ConnectedAt   string          `fesl:"curTime"`
	MessengerIP   string          `fesl:"messengerIp"`
	MessengerPort int             `fesl:"messengerPort"`
	TheaterIP     string          `fesl:"theaterIp"`
	TheaterPort   int             `fesl:"theaterPort"`
}

type domainPartition struct {
	Name    string `fesl:"domain"`
	SubName string `fesl:"subDomain"`
}

func (fm *FeslManager) hello(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	redisState := fm.createState(fmt.Sprintf(
		"%s-%s",
		event.Command.Message["clientType"],
		event.Client.IpAddr.String(),
	))

	event.Client.HashState = redisState

	if !fm.server {
		fm.gsumGetSessionID(event)
	}

	saveRedis := map[string]interface{}{
		"SDKVersion":     event.Command.Message["SDKVersion"],
		"clientPlatform": event.Command.Message["clientPlatform"],
		"clientString":   event.Command.Message["clientString"],
		"clientType":     event.Command.Message["clientType"],
		"clientVersion":  event.Command.Message["clientVersion"],
		"locale":         event.Command.Message["locale"],
		"sku":            event.Command.Message["sku"],
	}
	event.Client.HashState.SetM(saveRedis)

	ans := ansHello{
		Taxon:       fsysHello,
		ConnTTL:     int((1 * time.Hour).Seconds()),
		ConnectedAt: time.Now().Format("Jan-02-2006 15:04:05 MST"),
		TheaterIP:   config.General.ThtrAddr,
		//MessengerIP:   config.General.TelemetricsIP,
		//MessengerPort: config.General.TelemetricsPort,
	}

	if fm.server {
		ans.Domain = domainPartition{"eagames", "bfwest-server"}
		ans.TheaterPort = config.General.ThtrServerPort
	} else {
		ans.Domain = domainPartition{"eagames", "bfwest-dedicated"}
		ans.TheaterPort = config.General.ThtrClientPort
	}

	event.Client.WriteEncode(&codec.Packet{
		Payload: ans,
		Type:    fsys,
		Step:    0xC0000001,
	})
}

const (
	pingEUCentral = "iad"
)

type ansGetPingSites struct {
	Taxon     string     `fesl:"TXN"`
	MinPings  int        `fesl:"minPingSitesToPing"`
	PingSites []pingSite `fesl:"pingSites"`
}

type pingSite struct {
	Addr string `fesl:"addr"`
	Name string `fesl:"name"`
	Type int    `fesl:"type"`
}

// GetPingSites - returns a list of endpoints to test for the lowest latency on a client
func (fm *FeslManager) GetPingSites(event network.EventClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	event.Client.WriteEncode(&codec.Packet{
		Type: event.Command.Query,
		Step: event.Command.PayloadID,
		Payload: ansGetPingSites{
			Taxon:    fsysGetPingSites,
			MinPings: 1,
			PingSites: []pingSite{
				{"8.8.8.8", pingEUCentral, 0},
			},
		},
	})
}
