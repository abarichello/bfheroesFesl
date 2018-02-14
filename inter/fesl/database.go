package fesl

import (
	"database/sql"
	"strings"

	"github.com/sirupsen/logrus"
)

type Database struct {
	db   *sql.DB
	name string

	// Database Statements
	stmtGetUserByGameToken          *sql.Stmt
	stmtGetServerBySecret           *sql.Stmt
	stmtGetServerByID               *sql.Stmt
	stmtGetServerByName             *sql.Stmt
	stmtGetHeroesByUserID           *sql.Stmt
	stmtGetHeroeByName              *sql.Stmt
	stmtGetHeroeByID                *sql.Stmt
	stmtClearGameServerStats        *sql.Stmt
	mapGetStatsVariableAmount       map[int]*sql.Stmt
	mapGetServerStatsVariableAmount map[int]*sql.Stmt
	mapSetStatsVariableAmount       map[int]*sql.Stmt
	mapSetServerStatsVariableAmount map[int]*sql.Stmt
}

func NewDatabase(conn *sql.DB) (*Database, error) {
	db := &Database{db: conn}

	db.mapGetStatsVariableAmount = make(map[int]*sql.Stmt)
	db.mapGetServerStatsVariableAmount = make(map[int]*sql.Stmt)
	db.mapSetStatsVariableAmount = make(map[int]*sql.Stmt)

	// Prepare database statements
	db.prepareStatements()

	_, err := db.stmtClearGameServerStats.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) getServerStatsVariableAmount(statsAmount int) *sql.Stmt {
	var err error

	// Check if we already have a statement prepared for that amount of stats
	if statement, ok := d.mapGetServerStatsVariableAmount[statsAmount]; ok {
		return statement
	}

	var query string
	for i := 1; i < statsAmount; i++ {
		query += "?, "
	}

	sql := "SELECT gid, statsKey, statsValue" +
		"	FROM game_server_stats" +
		"	WHERE gid=?" +
		"		AND statsKey IN (" + query + "?)"

	d.mapGetServerStatsVariableAmount[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing mapGetServerStatsVariableAmount with "+sql+" query.", err.Error())
	}

	return d.mapGetServerStatsVariableAmount[statsAmount]
}

func (d *Database) getStatsStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check if we already have a statement prepared for that amount of stats
	if statement, ok := d.mapGetStatsVariableAmount[statsAmount]; ok {
		return statement
	}

	var query string
	for i := 1; i < statsAmount; i++ {
		query += "?, "
	}

	sql := "SELECT user_id, heroID, statsKey, statsValue" +
		"	FROM game_stats" +
		"	WHERE heroID=?" +
		"		AND user_id=?" +
		"		AND statsKey IN (" + query + "?)"

	d.mapGetStatsVariableAmount[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetStatsVariableAmount with "+sql+" query.", err.Error())
	}

	return d.mapGetStatsVariableAmount[statsAmount]
}

func (d *Database) setStatsStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check if we already have a statement prepared for that amount of stats
	if statement, ok := d.mapSetStatsVariableAmount[statsAmount]; ok {
		return statement
	}

	var query string
	for i := 1; i < statsAmount; i++ {
		query += "(?, ?, ?, ?), "
	}

	sql := "INSERT INTO game_stats" +
		"	(user_id, heroID, statsKey, statsValue)" +
		"	VALUES " + query + "(?, ?, ?, ?)" +
		"	ON DUPLICATE KEY UPDATE" +
		"	statsValue=VALUES(statsValue)"

	d.mapSetStatsVariableAmount[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing stmtSetStatsVariableAmount with "+sql+" query.", err.Error())
	}

	return d.mapSetStatsVariableAmount[statsAmount]
}

func (d *Database) prepareStatements() {
	var err error

	d.stmtGetUserByGameToken, err = d.db.Prepare(
		"SELECT id, username, email, birthday, language, country, game_token" +
			"	FROM users" +
			"	WHERE game_token = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetUserByGameToken.", err.Error())
	}

	d.stmtGetHeroesByUserID, err = d.db.Prepare(
		"SELECT id, user_id, heroName, online" +
			"	FROM game_heroes" +
			"	WHERE user_id = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroesByUserID.", err.Error())
	}

	d.stmtGetHeroeByName, err = d.db.Prepare(
		"SELECT id, user_id, heroName, online" +
			"	FROM game_heroes" +
			"	WHERE heroName = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroesByUserID.", err.Error())
	}

	d.stmtGetHeroeByID, err = d.db.Prepare(
		"SELECT id, user_id, heroName, online" +
			"	FROM game_heroes" +
			"	WHERE id = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroeByID.", err.Error())
	}

	d.stmtGetServerBySecret, err = d.db.Prepare(
		"SELECT game_servers.id, users.id, game_servers.servername, game_servers.secretKey, users.username" +
			"	FROM game_servers" +
			"	LEFT JOIN users" +
			"		ON users.id=game_servers.user_id" +
			"	WHERE secretKey = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetServerBySecret.", err.Error())
	}

	d.stmtGetServerByID, err = d.db.Prepare(
		"SELECT game_servers.id, users.id, game_servers.servername, game_servers.secretKey, users.username" +
			"	FROM game_servers" +
			"	LEFT JOIN users" +
			"		ON users.id=game_servers.user_id" +
			"	WHERE game_servers.id = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetServerByID.", err.Error())
	}

	d.stmtGetServerByName, err = d.db.Prepare(
		"SELECT game_servers.id, users.id, game_servers.servername, game_servers.secretKey, users.username" +
			"	FROM game_servers" +
			"	LEFT JOIN users" +
			"		ON users.id=game_servers.user_id" +
			"	WHERE game_servers.servername = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetServerByName.", err.Error())
	}

	d.stmtClearGameServerStats, err = d.db.Prepare(
		"DELETE FROM game_server_stats")
	if err != nil {
		logrus.Fatalln("Error preparing stmtClearGameServerStats.", err.Error())
	}
}

func (d *Database) closeStatements() {
	d.stmtGetUserByGameToken.Close()
	d.stmtGetServerBySecret.Close()
	d.stmtGetServerByID.Close()
	d.stmtGetServerByName.Close()
	d.stmtGetHeroesByUserID.Close()
	d.stmtGetHeroeByName.Close()
	d.stmtClearGameServerStats.Close()

	// Close the dynamic lenght getStats statements
	for index := range d.mapGetStatsVariableAmount {
		d.mapGetStatsVariableAmount[index].Close()
	}

	// Close the dynamic lenght setStats statements
	for index := range d.mapSetStatsVariableAmount {
		d.mapSetStatsVariableAmount[index].Close()
	}
}

// MysqlRealEscapeString - you know
func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
