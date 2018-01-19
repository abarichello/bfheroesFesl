package fesl

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Synaxis/unstable/backend/config"
	"github.com/Synaxis/unstable/backend/inter/network"
	"github.com/Synaxis/unstable/backend/storage/level"

	"github.com/sirupsen/logrus"
)

// FeslManager - handles incoming and outgoing FESL data
type FeslManager struct {
	name   string
	db     *Database
	level  *level.Level
	socket *network.Socket
	server bool
}

// New creates and starts a new ClientManager
func New(name, bind string, cert config.Fixtures, server bool, conn *sql.DB, lvl *level.Level) *FeslManager {
	db, err := NewDatabase(conn)
	if err != nil {
		return nil
	}

	socket, err := network.NewSocketTLS(name, bind, cert.Path, cert.PrivateKey)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	fm := &FeslManager{
		db:     db,
		level:  lvl,
		name:   name,
		server: server,
		socket: socket,
	}

	go fm.run()
	return fm
}

func (fm *FeslManager) run() {
	// Close all database statements
	defer fm.db.closeStatements()

	for {
		select {
		case event := <-fm.socket.EventChan:
			switch event.Name {
			case "newClient":
				fm.newClient(event.Data.(network.EventNewClient)) // TLS
			case "client.command.Hello":
				fm.hello(event.Data.(network.EventClientCommand))
			case "client.command.NuLogin":
				fm.NuLogin(event.Data.(network.EventClientCommand))
			case "client.command.NuGetPersonas":
				fm.NuGetPersonas(event.Data.(network.EventClientCommand))
			case "client.command.NuGetAccount":
				fm.NuGetAccount(event.Data.(network.EventClientCommand))
			case "client.command.NuLoginPersona":
				fm.NuLoginPersona(event.Data.(network.EventClientCommand))
			case "client.command.GetStatsForOwners":
				fm.GetStatsForOwners(event.Data.(network.EventClientCommand))
			case "client.command.GetStats":
				fm.GetStats(event.Data.(network.EventClientCommand))
			case "client.command.NuLookupUserInfo":
				fm.NuLookupUserInfo(event.Data.(network.EventClientCommand))
			case "client.command.GetPingSites":
				fm.GetPingSites(event.Data.(network.EventClientCommand))
			case "client.command.UpdateStats":
				fm.UpdateStats(event.Data.(network.EventClientCommand))
			case "client.command.GetTelemetryToken":
				fm.GetTelemetryToken(event.Data.(network.EventClientCommand))
			case "client.command.Start":
				fm.Start(event.Data.(network.EventClientCommand))
			case "client.close":
				fm.close(event.Data.(network.EventClientClose)) // TLS
			case "client.command":
				txn := event.Data.(network.EventClientCommand).Command.Message["TXN"]
				logrus.WithFields(logrus.Fields{
					"srv": fm.name,
					"cmd": fmt.Sprintf("%s/TXN:%s", event.Name, txn),
				}).Debugf("Got event")
			default:
				logrus.WithFields(logrus.Fields{"srv": fm.name, "event": event.Name}).Debugf("Got event")
			}
		}
	}
}

// TLS
func (fm *FeslManager) newClient(event network.EventNewClient) {
	if !event.Client.IsActive {
		logrus.Println("Client Left")
		return
	}

	fm.fsysMemCheck(&event)

	// Start Heartbeat
	event.Client.State.HeartTicker = time.NewTicker(time.Second * 10)
	go func() {
		for {
			if !event.Client.IsActive {
				return
			}
			select {
			case <-event.Client.State.HeartTicker.C:
				if !event.Client.IsActive {
					return
				}
				fm.fsysMemCheck(&event)
			}
		}
	}()

	logrus.Println("Client Connecting")

}

// TLS
func (fm *FeslManager) close(event network.EventClientClose) {
	logrus.Println("Client closed.")

	if event.Client.HashState != nil {
		if event.Client.HashState.Get("lkeys") != "" {
			lkeys := strings.Split(event.Client.HashState.Get("lkeys"), ";")
			for _, lkey := range lkeys {
				lkeyRedis := fm.level.NewObject("lkeys", lkey)
				lkeyRedis.Delete()
			}
		}

		event.Client.HashState.Delete()
	}

	if !event.Client.State.HasLogin {
		return
	}
}

func (fm *FeslManager) createState(ident string) *level.State {
	return fm.level.NewState(ident)
}
