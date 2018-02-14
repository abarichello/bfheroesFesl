package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var (
	General  cfg
	Database MySQL
	Cert     Fixtures
)

type cfg struct {
	LogLevel string         `envconfig:"LOG_LEVEL" default:"DEBUG"`

	HTTPBind  string `envconfig:"HTTP_BIND" default:"0.0.0.0:8080"`
	HTTPSBind string `envconfig:"HTTPS_BIND" default:"0.0.0.0:443"`


    ThtrAddr string          `envconfig:"THEATER_ADDR" default:"127.0.0.1"`    
	//TelemetricsIP   string `envconfig:"TELEMETRICS_IP" default:"127.0.0.1"`
	//TelemetricsPort int      `envconfig:"TELEMETRICS_PORT" default:"13505"`


	GameSpyIP    string       `envconfig:"GAMESPY_IP" default:"0.0.0.0"`
	FeslClientPort int    `envconfig:"FESL_CLIENT_PORT" default:"18270"`
	FeslServerPort int    `envconfig:"FESL_SERVER_PORT" default:"18051"`
	ThtrClientPort int `envconfig:"THEATER_CLIENT_PORT" default:"18275"`
	ThtrServerPort int `envconfig:"THEATER_SERVER_PORT" default:"18056"`



	LevelDBPath string `envconfig:"LEVEL_DB_PATH" default:"_data/lvl.db"`
}

type MySQL struct {
	UserName string  `envconfig:"DATABASE_USERNAME" default:"root"`
	Password string                 `envconfig:"DATABASE_PASSWORD"`
	Host     string `envconfig:"DATABASE_HOST" default:"127.0.0.1"`
	Port     int         `envconfig:"DATABASE_PORT" default:"3306"`
	Name     string     `envconfig:"DATABASE_NAME" default:"naomi"`
}

type Fixtures struct {
	Path            string `envconfig:"CERT_PATH" default:"./config/cert.pem"`
	PrivateKey string `envconfig:"PRIVATE_KEY_PATH" default:"./config/key.pem"`
}

func Initialize() {
	if err := envconfig.Process("", &General); err != nil {
		log.Fatal(err)
	}
	if err := envconfig.Process("", &Database); err != nil {
		log.Fatal(err)
	}
	if err := envconfig.Process("", &Cert); err != nil {
		log.Fatal(err)
	}
}

// LogLevel parses a default log level from a string
func LogLevel() logrus.Level {
	lvl, err := logrus.ParseLevel(General.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	return lvl
}

func bindAddr(addr string, port int) string {
	return fmt.Sprintf("%s:%d", addr, port)
}

func FeslClientAddr() string {
	return bindAddr(General.GameSpyIP, General.FeslClientPort)
}

func FeslServerAddr() string {
	return bindAddr(General.GameSpyIP, General.FeslServerPort)
}

func ThtrClientAddr() string {
	return bindAddr(General.GameSpyIP, General.ThtrClientPort)
}

func ThtrServerAddr() string {
	return bindAddr(General.GameSpyIP, General.ThtrServerPort)
}
