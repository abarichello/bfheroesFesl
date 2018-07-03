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
	stmtGetBookmark				*sql.Stmt
	stmtGetHeroByToken          *sql.Stmt
	stmtGetServerBySecret       *sql.Stmt
	stmtGetServerByID           *sql.Stmt
	stmtGetServerByName         *sql.Stmt
	stmtGetHeroesByUserID       *sql.Stmt
	stmtGetHeroByName           *sql.Stmt
	stmtGetHeroByID             *sql.Stmt
	stmtClearServerStats        *sql.Stmt
	MapGetStatsQuery       		map[int]*sql.Stmt
	MapGetServerStatsQuery 		map[int]*sql.Stmt
	MapSetStatsQuery       		map[int]*sql.Stmt
	MapSetServerStatsQuery 		map[int]*sql.Stmt
	// MapGetBookmark				map[int]*sql.Stmt	
}


func NewDatabase(conn *sql.DB) (*Database, error) {
	db := &Database{db: conn}

	db.MapGetStatsQuery = make(map[int]*sql.Stmt)
	db.MapGetServerStatsQuery = make(map[int]*sql.Stmt)
	db.MapSetStatsQuery = make(map[int]*sql.Stmt)

	// Prepare database statements
	db.prepareStatements()

	_, err := db.stmtClearServerStats.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}


// MysqlRealEscapeString - you know
func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

// func (d *Database) getBookmark(statsAmount int) *sql.Stmt {
// 	var err error

// 	// Check if Statement is prepared
// 	if statement, ok := d.MapGetServerStatsQuery[statsAmount]; ok {
// 		return statement
// 	}

// 	var query string

// 	for i := 1; i < statsAmount; i++ {
// 		query += "?, "
// 	}

// 	sql := "SELECT gid" +
// 		"	FROM game_server_player_preferences" +
// 		"	WHERE userid=?" +
// 		"		AND statsKey IN (" + query + "?)"

// 	d.MapGetServerStatsQuery[statsAmount], err = d.db.Prepare(sql)
// 	if err != nil {
// 		logrus.Println("Error preparing MapGetServerStatsQuery with "+sql+" query.", err.Error())
// 	}

// 	return d.MapGetBookmark[statsAmount]
// }


func (d *Database) getServerStatsQuery(statsAmount int) *sql.Stmt {
	var err error

	// Check if Statement is prepared
	if statement, ok := d.MapGetServerStatsQuery[statsAmount]; ok {
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

	d.MapGetServerStatsQuery[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing MapGetServerStatsQuery with "+sql+" query.", err.Error())
	}

	return d.MapGetServerStatsQuery[statsAmount]
}


func (d *Database) getStatsStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check if we have a statement prepared for the stats
	if statement, ok := d.MapGetStatsQuery[statsAmount]; ok {
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

	d.MapGetStatsQuery[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetStatsVariableAmount with "+sql+" query.", err.Error())
	}

	return d.MapGetStatsQuery[statsAmount]
}


func (d *Database) setStatsStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check if we have a statement prepared for the stats
	if statement, ok := d.MapSetStatsQuery[statsAmount]; ok {
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

	d.MapSetStatsQuery[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing stmtSetStatsVariableAmount with "+sql+" query.", err.Error())
	}

	return d.MapSetStatsQuery[statsAmount]
}


func (d *Database) prepareStatements() {
	var err error

	//Client Login/sessionID
	d.stmtGetHeroByToken, err = d.db.Prepare(
		"SELECT id, username, email, birthday, language, country, game_token" +
			"	FROM users" +
			"	WHERE game_token = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroByToken.", err.Error())
	}

	//Client Login
	d.stmtGetHeroesByUserID, err = d.db.Prepare(
		"SELECT id, user_id, heroName, online" +
			"	FROM game_heroes" +
			"	WHERE user_id = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroesByUserID.", err.Error())
	}

	//Client Login
	d.stmtGetHeroByName, err = d.db.Prepare(
		"SELECT id, user_id, heroName, online" +
			"	FROM game_heroes" +
			"	WHERE heroName = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroesByUserID.", err.Error())
	}
	//Client Login
	d.stmtGetHeroByID, err = d.db.Prepare(
		"SELECT id, user_id, heroName, online" +
			"	FROM game_heroes" +
			"	WHERE id = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroByID.", err.Error())
	}

	d.stmtGetBookmark, err = d.db.Prepare(
		"SELECT gid" +
		"	FROM game_player_server_preferences" +
		"	WHERE userid = ?")
		if err != nil {
			logrus.Println("Error Bookmark", err.Error())
		}

	//Server Login
	d.stmtGetServerBySecret, err = d.db.Prepare(
		"SELECT game_servers.id, users.id, game_servers.servername, game_servers.secretKey, users.username" +
			"	FROM game_servers" +
			"	LEFT JOIN users" +
			"		ON users.id=game_servers.user_id" +
			"	WHERE secretKey = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetServerBySecret.", err.Error())
	}

	//Server Login
	d.stmtGetServerByID, err = d.db.Prepare(
		"SELECT game_servers.id, users.id, game_servers.servername, game_servers.secretKey, users.username" +
			"	FROM game_servers" +
			"	LEFT JOIN users" +
			"		ON users.id=game_servers.user_id" +
			"	WHERE game_servers.id = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetServerByID.", err.Error())
	}

	//Server Login
	d.stmtGetServerByName, err = d.db.Prepare(
		"SELECT game_servers.id, users.id, game_servers.servername, game_servers.secretKey, users.username" +
			"	FROM game_servers" +
			"	LEFT JOIN users" +
			"		ON users.id=game_servers.user_id" +
			"	WHERE game_servers.servername = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetServerByName.", err.Error())
	}

	d.stmtClearServerStats, err = d.db.Prepare(
		"DELETE FROM game_server_stats")
	if err != nil {
		logrus.Fatalln("Error preparing stmtClearServerStats.", err.Error())
	}
}

func (d *Database) closeStatements() {
	d.stmtGetHeroByToken.Close()
	d.stmtGetServerBySecret.Close()
	d.stmtGetServerByID.Close()
	d.stmtGetServerByName.Close()
	d.stmtGetHeroesByUserID.Close()
	d.stmtGetHeroByName.Close()
	d.stmtClearServerStats.Close()

	// Close the dynamic lenght getStats statements
	for index := range d.MapGetStatsQuery {
		d.MapGetStatsQuery[index].Close()
	}

	// Close the dynamic lenght setStats statements
	for index := range d.MapSetStatsQuery {
		d.MapSetStatsQuery[index].Close()
	}
}

