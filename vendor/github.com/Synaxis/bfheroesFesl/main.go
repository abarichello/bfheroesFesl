package main

import (
	"database/sql"
	"flag"

	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/Synaxis/bfheroesFesl/inter/fesl"
	"github.com/Synaxis/bfheroesFesl/inter/theater"
	"github.com/Synaxis/bfheroesFesl/server"
	"github.com/Synaxis/bfheroesFesl/storage/database"
	"github.com/Synaxis/bfheroesFesl/storage/level"

	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

var (
	configFile string
)

func main() {
	initConfig()
	initLogger()

	mdb, _ := newMySQL()
	ldb, _ := newLevelDB()

	startServer(mdb, ldb)

	logrus.Println("Serving..")
	a := make(chan bool)
	<-a
}

func initConfig() {
	flag.StringVar(&configFile, "config", ".env", "Path to configuration file")

	flag.Parse()
	gotenv.Load(configFile)
	config.Initialize()
}

func initLogger() {
	logrus.SetLevel(config.LogLevel())
}

func newMySQL() (*sql.DB, error) {
	// DB Connection
	db, err := database.New(config.Database)
	if err != nil {
		logrus.Fatal("Error connecting to DB:", err)
	}
	return db, err
}

func newLevelDB() (*level.Level, error) {
	lvl, err := level.New(config.General.LevelDBPath, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	return lvl, err
}

func startServer(mdb *sql.DB, ldb *level.Level) {
	fesl.New("FM", config.FeslClientAddr(), config.Cert, false, mdb, ldb)
	fesl.New("SFM", config.FeslServerAddr(), config.Cert, true, mdb, ldb)

	theater.New("TM", config.ThtrClientAddr(), mdb, ldb)
	theater.New("STM", config.ThtrServerAddr(), mdb, ldb)

	srv := server.New(config.Cert)
	srv.ListenAndServe(
		config.General.HTTPBind,
		config.General.HTTPSBind,
	)
}
