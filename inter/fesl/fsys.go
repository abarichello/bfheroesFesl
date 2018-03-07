package fesl

import (
	"fmt"
	"time"

	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

//THIS IS THE HELLO Pkt ;)
const (
	fsys             = "fsys"
	fsysGetPingSites = "GetPingSites"
	fsysHello        = "Hello"
	fsysMemCheck     = "MemCheck"
	fsysPing         = "Ping"
)

type ansMemCheck struct {
	TXN       string     `fesl:"TXN"`
	MemChecks []memCheck `fesl:"memcheck"`
	Salt      string     `fesl:"salt"`
}

type memCheck struct {
	Length int    `fesl:"len"`
	Addr   string `fesl:"addr"`
}

func (fm *FeslManager) fsysMemCheck(event *network.EventNewClient) {
	event.Client.Answer(&codec.Pkt{
		Type: fsys,
		Content: ansMemCheck{
			TXN:  fsysMemCheck,
			Salt: "5",
		},
		Send: 0xC0000000,
	})
}

type ansHello struct {
	TXN           string          `fesl:"TXN"`
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

func (fm *FeslManager) hello(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	redisState := fm.createState(fmt.Sprintf(
		"%s-%s",
		event.Process.Msg["clientType"],
		event.Client.IpAddr.String(),
	))

	event.Client.HashState = redisState

	if !fm.server {
		fm.gsumGetSessionID(event)
	}

	saveRedis := map[string]interface{}{
		"SDKVersion":     event.Process.Msg["SDKVersion"],
		"clientPlatform": event.Process.Msg["clientPlatform"],
		"clientString":   event.Process.Msg["clientString"],
		"clientType":     event.Process.Msg["clientType"],
		"clientVersion":  event.Process.Msg["clientVersion"],
		"locale":         event.Process.Msg["locale"],
		"sku":            event.Process.Msg["sku"],
	}
	event.Client.HashState.SetM(saveRedis)

	ans := ansHello{
		TXN:         fsysHello,
		ConnTTL:     int((1 * time.Hour).Seconds()),
		ConnectedAt: time.Now().Format("Jan-02-2006 15:04:05 MST"),
		TheaterIP:   config.General.ThtrAddr,
	}

	if fm.server {
		ans.Domain = domainPartition{"eagames", "bfwest-server"}
		ans.TheaterPort = config.General.ThtrServerPort
	} else {
		ans.Domain = domainPartition{"eagames", "bfwest-dedicated"}
		ans.TheaterPort = config.General.ThtrClientPort
	}

	event.Client.Answer(&codec.Pkt{
		Content: ans,
		Type:    fsys,
		Send:    0xC0000001,
	})
}

const (
	location = "iad"
)

//is this usefull ? added google dns to test

type ansGetPingSites struct {
	TXN       string     `fesl:"TXN"`
	MinPings  int        `fesl:"minPingSitesToPing"`
	PingSites []pingSite `fesl:"pingSites"`
}

type pingSite struct {
	Addr string `fesl:"addr"`
	Name string `fesl:"name"`
	Type int    `fesl:"type"`
}

// GetPingSites - Get Pings for something
func (fm *FeslManager) GetPingSites(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	event.Client.Answer(&codec.Pkt{
		Type: event.Process.Query,
		Send: event.Process.HEX,
		Content: ansGetPingSites{
			TXN:      fsysGetPingSites,
			MinPings: 1,
			PingSites: []pingSite{
				{"8.8.8.8",
					location,
					0}, //0 = MS
			},
		},
	})
}
