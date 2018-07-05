package theater

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type Database struct {
	db *sql.DB

	// Database Statements
	stmtGetHeroByID            *sql.Stmt
	stmtDeleteServerStatsByGID *sql.Stmt
	stmtDeleteGameByGID        *sql.Stmt
	stmtAddGame                *sql.Stmt
	stmtGameIncreaseJoining    *sql.Stmt
	UpdateGid                  *sql.Stmt
	stmtGameIncreaseTeam1      *sql.Stmt
	stmtGameIncreaseTeam2      *sql.Stmt
	stmtGameDecreaseTeam1      *sql.Stmt
	stmtGameDecreaseTeam2      *sql.Stmt
	stmtUpdateGame             *sql.Stmt
	stmtCreateServer           *sql.Stmt
	MapGetStatsQuery           map[int]*sql.Stmt
	MapSetServerStatsQuery     map[int]*sql.Stmt
	MapSetPlayerStatsQuery     map[int]*sql.Stmt
}

func NewDatabase(conn *sql.DB) (*Database, error) {
	db := &Database{db: conn}

	// Prepare DB statements
	db.MapGetStatsQuery = make(map[int]*sql.Stmt)
	db.MapSetServerStatsQuery = make(map[int]*sql.Stmt)
	db.MapSetPlayerStatsQuery = make(map[int]*sql.Stmt)
	db.prepareStatements()

	return db, nil
}

func (d *Database) prepareStatements() {
	var err error

	d.stmtGetHeroByID, err = d.db.Prepare(
		"SELECT id, user_id, heroName, online" +
			"	FROM game_heroes" +
			"	WHERE id = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetHeroByID.", err.Error())
	}

	d.stmtDeleteServerStatsByGID, err = d.db.Prepare(
		"DELETE FROM game_server_stats WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtClearServerStats.", err.Error())
	}

	d.stmtDeleteGameByGID, err = d.db.Prepare(
		"DELETE FROM games WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtClearServerStats.", err.Error())
	}

	d.stmtAddGame, err = d.db.Prepare(
		"INSERT INTO games (" +
			"	gid," +
			"	game_ip," +
			"	game_port," +
			"	game_version," +
			"	status_join," +
			"	status_mapname," +
			"	players_connected," +
			"	players_joining," +
			"	players_max," +
			"	team_1," +
			"	team_2," +
			"	team_distribution," +
			"	created_at," +
			"	updated_at)" +
			"VALUES" +
			"	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())")
	if err != nil {
		logrus.Fatalln("Error preparing stmtAddGame:", err.Error())
	}

	d.stmtGameIncreaseJoining, err = d.db.Prepare(
		"UPDATE games SET " +
			"	players_joining = players_joining + 1," +
			"	updated_at = NOW()" +
			"WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGameIncreaseJoining.", err.Error())
	}

	d.stmtGameIncreaseTeam1, err = d.db.Prepare(
		"UPDATE games SET " +
			"	players_connected = players_connected + 1," +
			"	players_joining = players_joining - 1," +
			"	team_1 = team_1 + 1," +
			"	updated_at = NOW()" +
			"WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGameIncreaseTeam1.", err.Error())
	}

	d.stmtGameIncreaseTeam2, err = d.db.Prepare(
		"UPDATE games SET " +
			"	players_connected = players_connected + 1," +
			"	players_joining = players_joining - 1," +
			"	team_2 = team_2 + 1," +
			"	updated_at = NOW()" +
			"WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGameIncreaseTeam2.", err.Error())
	}

	d.stmtGameDecreaseTeam1, err = d.db.Prepare(
		"UPDATE games SET " +
			"	players_connected = players_connected - 1," +
			"	team_1 = team_1 - 1," +
			"	updated_at = NOW()" +
			"WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGameDecreaseTeam1.", err.Error())
	}

	d.stmtGameDecreaseTeam2, err = d.db.Prepare(
		"UPDATE games SET " +
			"	players_connected = players_connected - 1," +
			"	team_2 = team_2 - 1," +
			"	updated_at = NOW()" +
			"WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtGameDecreaseTeam2.", err.Error())
	}

	d.stmtUpdateGame, err = d.db.Prepare(
		"UPDATE games SET" +
			"	updated_at = NOW()" +
			"WHERE gid = ?")
	if err != nil {
		logrus.Fatalln("Error preparing stmtUpdateGame.", err.Error())
	}

	d.stmtCreateServer, err = d.db.Prepare(
		"INSERT INTO game_server_client (" +
			"name, " +
			"community_name, " +
			"ip_address, " +
			"port, " +
			"client_version) " +
			"VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		logrus.Fatalln("Error preparing stmtCreateServer.", err.Error())
	}
}
func (d *Database) getStatsStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check Statement is OK for query
	if statement, ok := d.MapGetStatsQuery[statsAmount]; ok {
		return statement
	}

	var query string
	for i := 1; i < statsAmount; i++ {
		query += "?, "
	}

	sql := "SELECT game_heroes.user_id, game_heroes.id, game_heroes.heroName, game_stats.statsKey, game_stats.statsValue" +
		"	FROM game_heroes" +
		"	LEFT JOIN game_stats" +
		"		ON game_stats.user_id = game_heroes.user_id" +
		"		AND game_stats.heroID = game_heroes.id" +
		"	WHERE game_heroes.id=?" +
		"		AND game_stats.statsKey IN (" + query + "?)"

	d.MapGetStatsQuery[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing stmtGetStatsVariableAmount with "+sql+" query.", err.Error())
	}

	return d.MapGetStatsQuery[statsAmount]
}

func (d *Database) setServerStatsStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check Statement is OK for query
	if statement, ok := d.MapSetServerStatsQuery[statsAmount]; ok {
		return statement
	}

	var query string
	for i := 1; i < statsAmount; i++ {
		query += "(?, ?, ?, NOW()), "
	}

	sql := "INSERT INTO game_server_stats" +
		"	(gid, statsKey, statsValue, created_at)" +
		"	VALUES " + query + "(?, ?, ?, NOW())" +
		"	ON DUPLICATE KEY UPDATE" +
		"	statsValue=VALUES(statsValue)," +
		"   updated_at=NOW()"

	d.MapSetServerStatsQuery[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing setServerStatsStatement with "+sql+" query.", err.Error())
	}

	return d.MapSetServerStatsQuery[statsAmount]
}

func (d *Database) setServerGidStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check Statement is OK for query
	if statement, ok := d.MapSetServerStatsQuery[statsAmount]; ok {
		return statement
	}

	var query string
	for i := 1; i < statsAmount; i++ {
		query += "(?, ?, ?, NOW()), "
	}

	sql := "INSERT INTO game_server_stats" +
		"	(gid)" +
		"	VALUES " + query + "(?)" +
		"	ON DUPLICATE KEY UPDATE"

	d.MapSetServerStatsQuery[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing setServerStatsStatement with "+sql+" query.", err.Error())
	}

	return d.MapSetServerStatsQuery[statsAmount]
}

func (d *Database) setServerPlayerStatsStatement(statsAmount int) *sql.Stmt {
	var err error

	// Check Statement is OK for query
	if statement, ok := d.MapSetPlayerStatsQuery[statsAmount]; ok {
		return statement
	}

	var query string
	for i := 1; i < statsAmount; i++ {
		query += "(?, ?, ?, ?, NOW()), "
	}

	sql := "INSERT INTO game_server_player_stats" +
		"	(gid, pid, statsKey, statsValue, created_at)" +
		"	VALUES " + query + "(?, ?, ?, ?, NOW())" +
		"	ON DUPLICATE KEY UPDATE" +
		"	statsValue=VALUES(statsValue)," +
		"   updated_at=NOW()"

	d.MapSetPlayerStatsQuery[statsAmount], err = d.db.Prepare(sql)
	if err != nil {
		logrus.Fatalln("Error preparing MapSetPlayerStatsQuery with "+sql+" query.", err.Error())
	}

	return d.MapSetPlayerStatsQuery[statsAmount]
}

func (d *Database) closeStatements() {
	// Close the dynamic lenght getStats statements
	for index := range d.MapGetStatsQuery {
		d.MapGetStatsQuery[index].Close()
	}
}
