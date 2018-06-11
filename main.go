package main

import (
	"database/sql"
	"flag"
	"fmt"

	"strconv"
	"strings"

	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/Synaxis/bfheroesFesl/inter/fesl"
	"github.com/Synaxis/bfheroesFesl/inter/theater"
	"github.com/Synaxis/bfheroesFesl/storage/database"
	"github.com/Synaxis/bfheroesFesl/storage/level"

	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

func main() {
	initConfig()
	initLogger()

	mdb, _ := newMySQL()
	ldb, _ := newLevelDB()

	startServer(mdb, ldb)

	logrus.Println(" Fesl    Online")
    logrus.Println(" Theater Online")

	a := make(chan bool)
	<-a
}

func initConfig() {
	var (
		configFile string
	)
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

// separate tags to be recognized as key-value
func separateTags(singleTag string) (string, string) {
	tuple := strings.SplitN(singleTag, ":", 2)
	return tuple[0], tuple[1]
}

type FieldTag struct {
	tags map[string]string
}

func (ft *FieldTag) StringVal(tag string) (string, error) {
	value, ok := ft.tags[tag]
	if !ok {
		return "", fmt.Errorf("tag: %s not found", tag)
	}

	removedQuotes := strings.Trim(value, `"`)

	return removedQuotes, nil
}


func (ft *FieldTag) StringArr(tag string) ([]string, error) {
	s, err := ft.StringVal(tag)
	if err != nil {
		return nil, err
	}

	return strings.Split(s, ","), nil
}

func (ft *FieldTag) IntVal(tag string) (int, error) {
	s, err := ft.StringVal(tag)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func startServer(mdb *sql.DB, ldb *level.Level) {
	fesl.New("FM", config.FeslClientAddr(), false, mdb, ldb)
	fesl.New("SFM", config.FeslServerAddr(), true, mdb, ldb)

	theater.New("TM", config.ThtrClientAddr(), mdb, ldb)
	theater.New("STM", config.ThtrServerAddr(), mdb, ldb)
}



