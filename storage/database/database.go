package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Synaxis/bfheroesFesl/config"

	// Needed since we are using this as driver for MySQL database
	"github.com/Synaxis/logrus"
	_ "github.com/Synaxis/mysql"
)

// New tries to establish connection with database
func New(cfg config.MySQL) (*sql.DB, error) {
	return newMySQL(mysqlDSN(cfg))
}

// newMySQL establishes connection to MySQL database, and then pings it to
// verify if it responds
func newMySQL(dsnAddr string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", dsnAddr)
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

// mysqlDSN creates a DSN string to connect local instance of MySQL database
func mysqlDSN(cfg config.MySQL) string {
	connParams := url.Values{
		"charset":   {"utf8"},
		"parseTime": {"True"},
		"loc":       {"UTC"},
	}

	dsnAddr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		connParams.Encode(),
	)

	return dsnAddr
}
