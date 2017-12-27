package main

import (
	"database/sql"
	"flag"

	"bitbucket.org/openheroes/backend/config"
	"bitbucket.org/openheroes/backend/internal/fesl"
	"bitbucket.org/openheroes/backend/internal/theater"
	"bitbucket.org/openheroes/backend/server"
	"bitbucket.org/openheroes/backend/storage/database"
	"bitbucket.org/openheroes/backend/storage/level"

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

	// Use "github.com/google/gops/agent" to analyze resources
	// if err := agent.Listen(&agent.Options{}); err != nil {
	// 	log.Fatal(err)
	// }

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

	// logrus.SetFormatter(&logrus.JSONFormatter{
	// 	DisableTimestamp: true,
	// })
	// logrus.SetFormatter(new(prefixed.TextFormatter))
	// logrus.SetFormatter(&prefixed.TextFormatter{
	// 	DisableTimestamp: true,
	// 	DisableColors:    true,
	// })
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
	// Redis Connection
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: "192.168.33.10:6379",
	// })
	// if _, err = redisClient.Ping().Result(); err != nil {
	// 	log.Fatalln("Error connecting to redis:", err)
	// }

	// lvl, err := level.New("_data/lvl.db", redisClient)
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
