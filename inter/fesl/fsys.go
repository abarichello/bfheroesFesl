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

func (fm *Fesl) hello(event network.EvProcess) {
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

	answer := ansHello{
		TXN:         fsysHello,
		ConnTTL:     int((60 * time.Hour).Seconds()),
		ConnectedAt: time.Now().Format("Jan-02-2006 15:04:05 MST"),
		TheaterIP:   config.General.ThtrAddr,
		MessengerIP: config.General.MessengerAddr,
	}

	if fm.server {
		answer.Domain = domainPartition{"eagames", "bfwest-server"}
		answer.TheaterPort = config.General.ThtrServerPort
	} else {
		answer.Domain = domainPartition{"eagames", "bfwest-dedicated"}
		answer.TheaterPort = config.General.ThtrClientPort
	}	

	event.Client.Answer(&codec.Packet{
		Content: answer,
		Message: fsys,
		Send:    0xC0000001,
	})
}

type ansMemCheck struct {
	TXN      string `fesl:"TXN"`
	Salt     string `fesl:"salt"`
	memcheck string `fesl:"memcheck.[]"`
}

func (fm *Fesl) fsysMemCheck(event *network.EventNewClient) {	
	event.Client.Answer(&codec.Packet{
		Message: fsys,
		Content: ansMemCheck{
			TXN:      "MemCheck",
			memcheck: "0",
			Salt:     "0",
		},
		Send: 0xC0000000,
	})
}


///////////////////////////////////////////////
type ansGoodbye struct {
	TXN       string     `fesl:"TXN"`
	Reason    string     `fesl:"reason"`
	messageArr string    `fesl:"message"`
}

// Goodbye - Handle Client Close
func (fm *Fesl) Goodbye(event network.EvProcess) {	
	logrus.Println("Client Disconnected")	
	event.Client.Answer(&codec.Packet{
		Message: event.Process.Query,
		Send:    event.Process.HEX,
		Content: ansGoodbye{
			TXN:      "Goodbye",
			Reason:   "GOODBYE_CLIENT_NORMAL",
			messageArr: "n/a",			
			},
		},
	)
}

///////////////////////////////////////////////
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

// GetPingSites - Was used for Load-Balancer / Not working Now (but it's requested)
func (fm *Fesl) GetPingSites(event network.EvProcess) {
	if !event.Client.IsActive {
		return
	}

	event.Client.Answer(&codec.Packet{
		Message: event.Process.Query,
		Send:    event.Process.HEX,
		Content: ansGetPingSites{
			TXN:      fsysGetPingSites,
			MinPings: 1,
			PingSites: []pingSite{
				{"127.0.0.1", "iad", 1},
			},
		},
	})
}
