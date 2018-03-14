package fesl

import (
	"fmt"
	"time"

	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

//THIS IS THE HELLO PACKET
const (
	fsys             = "fsys"
	fsysGetPingSites = "GetPingSites"
	fsysHello        = "Hello"
	fsysPing         = "Ping"
)

type ansMemCheck struct {
	TXN      string `fesl:"TXN"`
	Salt     string `fesl:"salt"`
	mtype    string `fesl:"type"`
	memcheck string `fesl:"memcheck.[]`
	result   string `fesl:"result"`
}

func (fm *FeslManager) fsysMemCheck(event *network.EventNewClient) {
	event.Client.Answer(&codec.Packet{
		Message: fsys,
		Content: ansMemCheck{
			TXN:      "MemCheck",
			memcheck: "0",
			Salt:     "3",
			result:   "",
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
		MessengerIP: config.General.MessengerAddr,
	}

	if fm.server {
		ans.Domain = domainPartition{"eagames", "bfwest-server"}
		ans.TheaterPort = config.General.ThtrServerPort
	} else {
		ans.Domain = domainPartition{"eagames", "bfwest-dedicated"}
		ans.TheaterPort = config.General.ThtrClientPort
	}

	hex := event.Process.HEX
	event.Client.Answer(&codec.Packet{
		Content: ans,
		Message: fsys,
		Send:    hex,
	})
}

const (
	location = "iad"
)

type ansGetPingSites struct {
	TXN       string     `fesl:"TXN"`
	MinPings  int        `fesl:"minPingSitesToPing"`
	PingSites []pingSite `fesl:"pingSites"`
}

type pingSite struct {
	Addr    string `fesl:"addr"`
	Name    string `fesl:"name"`
	Message int    `fesl:"type"`
}

// GetPingSites - Get Pings for something
func (fm *FeslManager) GetPingSites(event network.EventClientProcess) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}

	hex := event.Process.HEX
	event.Client.Answer(&codec.Packet{
		Message: event.Process.Query,
		Send:    hex,
		Content: ansGetPingSites{
			TXN:      fsysGetPingSites,
			MinPings: 1,
			PingSites: []pingSite{
				{"127.0.0.1", location, 1},
			},
		},
	})
}
