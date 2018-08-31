package fesl

import (
	"database/sql"
	"fmt"
	"github.com/OSHeroes/bfheroesFesl/inter/network"
	"github.com/OSHeroes/bfheroesFesl/storage/level"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// TXN stands for Taxon, sub-query name of the command
// Fesl - handles incoming and outgoing FESL data
type Fesl struct {
	name   string
	db     *Database
	level  *level.Level
	socket *network.Socket

	server bool
}

// New creates and starts a new ClientManager
func New(name, bind string, server bool, conn *sql.DB, lvl *level.Level) *Fesl {
	db, err := NewDatabase(conn)
	if err != nil {
		return nil
	}

	socket, err := network.NewSocketTLS(name, bind)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	fm := &Fesl{
		db:     db,
		level:  lvl,
		name:   name,
		server: server,
		socket: socket,
	}

	go fm.run()
	return fm
}

func (fm *Fesl) run() {
	// Close all database statements
	defer fm.db.closeStatements()

	for {
		select {
		case event := <-fm.socket.EventChan:
			switch event.Name {
			case "newClient":
				fm.newClient(event.Data.(network.EventNewClient)) // TLS
			case "client.command.Hello":
				fm.hello(event.Data.(network.EvProcess))
			case "client.command.Telemetry":
				fm.Telemetry(event.Data.(network.EvProcess))
			case "client.command.NuLogin":
				fm.NuLogin(event.Data.(network.EvProcess))
			case "client.command.NuGetPersonas":
				fm.NuGetPersonas(event.Data.(network.EvProcess))
			case "client.command.NuGetPersonasServer":
				fm.NuGetPersonasServer(event.Data.(network.EvProcess))
			case "client.command.NuGetAccount":
				fm.NuGetAccount(event.Data.(network.EvProcess))
			case "client.command.GetStats":
				fm.GetStats(event.Data.(network.EvProcess))
			case "client.command.NuLookupUserInfo":
				fm.NuLookupUserInfo(event.Data.(network.EvProcess))
			case "client.command.NuLoginPersona":
				fm.NuLoginPersona(event.Data.(network.EvProcess))
			case "client.command.GetStatsForOwners":
				fm.GetStatsForOwners(event.Data.(network.EvProcess))
			case "client.command.GetPingSites":
				fm.GetPingSites(event.Data.(network.EvProcess))
			case "client.command.UpdateStats":
				fm.UpdateStats(event.Data.(network.EvProcess))
			case "client.command.Start":
				fm.Start(event.Data.(network.EvProcess))
			case "client.command.Goodbye":
				fm.Goodbye(event.Data.(network.EvProcess))
			case "client.close":
				fm.close(event.Data.(network.EventClientClose)) // TLS
			case "client.command":
				txn := event.Data.(network.EvProcess).Process.Msg["TXN"]
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
func (fm *Fesl) newClient(event network.EventNewClient) {
	fm.fsysMemCheck(&event)

	logrus.Println("Client Connecting")
	// Start Heartbeat
	event.Client.State.HeartTicker = time.NewTicker(time.Second * 10)
	go func() {
		for {
			if !event.Client.IsActive {
				return
			}
			select {
			case <-event.Client.State.HeartTicker.C:
				fm.fsysMemCheck(&event)
			}
		}
	}() //-> go routine

}

// TLS
func (fm *Fesl) close(event network.EventClientClose) {
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

func (fm *Fesl) createState(ident string) *level.State {
	return fm.level.NewState(ident)
}
