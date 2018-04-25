package theater

import (
	"database/sql"
	"time"

	"github.com/Synaxis/bfheroesFesl/inter/mm"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/storage/level"

	"github.com/sirupsen/logrus"
)

const (
	thtrCGAM = "CGAM"
	thtrCONN = "CONN"
	thtrECHO = "ECHO"
	thtrEGAM = "EGAM"
	thtrEGEG = "EGEG"
	thtrEGRQ = "EGRQ"
	thtrEGRS = "EGRS"
	thtrUQUE = "UQUE" //testing
	thtrENCL = "ENCL"
	thtrGDAT = "GDAT"
	thtrKICK = "KICK"
	thtrLDAT = "LDAT"
	thtrHTSN = "HTSN" //testing 
	thtrLLST = "LLST"
	thtrPENT = "PENT"
	thtrPING = "PING"
	thtrPLVT = "PLVT"
	thtrUBRA = "UBRA"
	thtrUPLA = "UPLA"
	thtrGREM = "GREM" //testing
	thtrUSER = "USER"
	thtrRGAM = "RGAM" //testing
)

// Theater Handles incoming and outgoing theater communication
type Theater struct {
	name      string
	socket    *network.Socket
	socketUDP *network.SocketUDP
	db        *Database
	level     *level.Level
}

// New creates and starts a new TheaterManager
func New(name, bind string, conn *sql.DB, lvl *level.Level) *Theater {
	socket, err := network.NewSocketTCP(name, bind, true)
	if err != nil {
		return nil
	}

	socketUDP, err := network.NewSocketUDP(name, bind, true)
	if err != nil {
		return nil
	}

	db, err := NewDatabase(conn)
	if err != nil {
		return nil
	}

	tm := &Theater{
		db:        db,
		level:     lvl,
		name:      name,
		socket:    socket,
		socketUDP: socketUDP,
	}
	go tm.Listen()

	return tm
}

func (tm *Theater) Listen() {
	defer tm.db.closeStatements()

	for {
		select {
		case event := <-tm.socketUDP.EventChan:
			switch event.Name {
			case "command.ECHO":
				go tm.ECHO(event)
			case "command":
				logrus.Debugf(" %s: %v", event.Name, event.Data.(*network.ProcessFESL))
			default:
				logrus.Debugf(" %s: %v", event.Name, event.Data)
			}
		case event := <-tm.socket.EventChan:
			switch event.Name {
			case "newClient":
				go tm.newClient(event.Data.(network.EventNewClient))
			case "client.command.CONN":
				go tm.CONN(event.Data.(network.EvProcess))
			case "client.command.USER":
				go tm.USER(event.Data.(network.EvProcess))
			case "client.command.GDAT":
				go tm.GDAT(event.Data.(network.EvProcess))
			case "client.command.EGAM":
				go tm.EGAM(event.Data.(network.EvProcess))
			case "client.command.ECNL":
				go tm.ECNL(event.Data.(network.EvProcess))
			case "client.command.CGAM":
				go tm.CGAM(event.Data.(network.EvProcess))
			case "client.command.UBRA":
				go tm.UBRA(event.Data.(network.EvProcess))
			case "client.command.RGAM":
				go tm.UGAM(event.Data.(network.EvProcess))
			case "client.command.UGAM":
				go tm.UGAM(event.Data.(network.EvProcess))
			case "client.command.EGRS":
				go tm.EGRS(event.Data.(network.EvProcess))
			case "client.command.UQUE":
				go tm.UQUE(event.Data.(network.EvProcess))
			case "client.command.PENT":
				go tm.PENT(event.Data.(network.EvProcess))
			case "client.command.PLVT":
				go tm.PLVT(event.Data.(network.EvProcess))
			case "client.command.UPLA":
				go tm.UPLA(event.Data.(network.EvProcess))
			case "client.command.HTSN":
				go tm.UPLA(event.Data.(network.EvProcess))
			 case "client.command.GREM":
			 	go tm.GREM(event.Data.(network.EvProcess))
			case "client.close":
				tm.close(event.Data.(network.EventClientClose))
			case "client.command":
				logrus.WithFields(logrus.Fields{" ": tm.name, "cmd": event.Name})
			default:
				logrus.WithFields(logrus.Fields{" ": tm.name, "cmd": event.Name})
			}
		}
	}
}

func (tM *Theater) error(event network.EventClientError) {
	logrus.Println("Client threw an error: ", event.Error)
}

func (tm *Theater) newClient(event network.EventNewClient) {
	if !event.Client.IsActive {
		logrus.Println("Cli Left")
		return
	}
	logrus.Println("Join Theather")

	// Start Heartbeat
	event.Client.State.HeartTicker = time.NewTicker(time.Second * 55)
	go func() {
		for event.Client.IsActive {
			select {
			case <-event.Client.State.HeartTicker.C:
				if !event.Client.IsActive {
					return
				}
			}
		}
	}()
}

func (tm *Theater) close(event network.EventClientClose) {
	logrus.Println("Client closed.")

	if event.Client.HashState != nil {

		if event.Client.HashState.Get("gdata:GID") != "" {

			// Delete game from db
			_, err := tm.db.stmtDeleteServerStatsByGID.Exec(event.Client.HashState.Get("gdata:GID"))
			if err != nil {
				logrus.Errorln("Failed deleting settings for  "+event.Client.HashState.Get("gdata:GID"), err.Error())
			}

			_, err = tm.db.stmtDeleteGameByGIDAnd.Exec(event.Client.HashState.Get("gdata:GID"))
			if err != nil {
				logrus.Errorln("Failed deleting game for "+event.Client.HashState.Get("gdata:GID"), err.Error())
			}

			// Delete game out of matchmaking array
			delete(mm.Games, event.Client.HashState.Get("gdata:GID"))

			gameServer := tm.level.NewObject("gdata", event.Client.HashState.Get("gdata:GID"))
			gameServer.Delete()
		}

		event.Client.HashState.Delete()
	}

	if !event.Client.State.HasLogin {
		return
	}
}
