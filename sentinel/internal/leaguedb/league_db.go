package leaguedb

import (
	"database/sql"
	"fmt"
	"log"
	"sentinel/internal/riotapi"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type LeagueDB struct {
	Config      map[string]string
	Multipliers map[string]any
	DB          *sql.DB
}

type Lobby struct {
	Id             string
	Creation_time  time.Time
	Creator_puuid  string
	Started        bool
	Start_time     *time.Time
	Matches        int
	Last_processed time.Time
}

type Player struct {
	Puuid       string
	Name        string
	Region      string
	Last_match  *time.Time
	Total_score int
	Completed   bool
}

func (leagueDB *LeagueDB) SetConfig() {
	// Read config file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	// Create a map of all configs and return the map
	configMap := make(map[string]string)
	for _, v := range viper.AllKeys() {
		configMap[strings.ToUpper(v)] = viper.Get(v).(string)
	}
	leagueDB.Config = configMap
}

// --- Bunch of helper functions ---

// Read multipliers for scoring.toml file located in root of project
func (leaguedb *LeagueDB) GetMultipliers() {
	viper.SetConfigFile("scoring.toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	multipliersMap := make(map[string]any)
	for _, v := range viper.AllKeys() {
		multipliersMap[v] = viper.Get(v)
	}
	leaguedb.Multipliers = multipliersMap
}
func (db *LeagueDB) ConnectToDB() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db.Config["MYSQL_USER"], db.Config["MYSQL_PASSWORD"],
		db.Config["MYSQL_URL"], db.Config["MYSQL_PORT"], db.Config["MYSQL_DATABASE"])

	mysql_db, err := sql.Open("mysql", dataSource)

	if err != nil {
		log.Fatalf("Error while connecting to MySQL DB %s", err)
	}

	db.DB = mysql_db
}

func (db *LeagueDB) GetLobbiesByLastProcessed() []Lobby {
	query := `
	SELECT id, creation_time, creator_puuid, started, started_time, matches, last_processed
	FROM lobbies
	WHERE started = TRUE
	ORDER BY last_processed
	LIMIT 10;`

	rows, err := db.DB.Query(query)

	if err != nil {
		log.Println(err)
	}

	var lobbies []Lobby

	for rows.Next() {

		lobby := Lobby{}

		rows.Scan(
			&lobby.Id, &lobby.Creation_time, &lobby.Creator_puuid, &lobby.Started,
			&lobby.Start_time, &lobby.Matches, &lobby.Last_processed,
		)
		lobbies = append(lobbies, lobby)
	}

	if lobbies == nil {
		log.Println("No lobbies found")
	}

	return lobbies
}

func (leagueDB *LeagueDB) GetPlayersInLobby(lobbyID string) []Player {

	players := []Player{}

	query := `
	SELECT puuid, name, region, last_match, total_score, completed
	FROM players
	WHERE lobby_id = ?`

	rows, err := leagueDB.DB.Query(query, lobbyID)

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		player := Player{}

		rows.Scan(&player.Puuid, &player.Name, &player.Region, &player.Last_match,
			&player.Total_score, &player.Completed)

		players = append(players, player)
	}

	if players == nil {
		log.Printf("No players found in lobby %s\n", lobbyID)
	}

	return players
}

func (leagueDB *LeagueDB) InsertMatchInfo(matchInfo riotapi.MatchInfo) {

	query := `
			INSERT INTO matches (match_id, player_puuid, champion, position, kills, deaths, assists, turrets, inhibs,
								dragons, rifts, barons, vision_score, creep_score, pentas, quadras, triples, doubles, win)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := leagueDB.DB.Exec(
		query, matchInfo.Match_ID, matchInfo.Puuid, matchInfo.Champion, matchInfo.Position, matchInfo.Kills, matchInfo.Deaths,
		matchInfo.Assists, matchInfo.Turrets, matchInfo.Inhibs, matchInfo.Dragons, matchInfo.Rifts, matchInfo.Barons, matchInfo.Vision_Score,
		matchInfo.Creep_Score, matchInfo.Pentas, matchInfo.Quadras, matchInfo.Triples, matchInfo.Doubles, matchInfo.Win)

	if err != nil {
		log.Println(err)
	}
}

func (leagueDB *LeagueDB) UpdateLastMatch(puuid string, last_match time.Time) {

	query := `
			UPDATE players
			SET last_match = ?
			WHERE puuid = ?`

	_, err := leagueDB.DB.Exec(query, last_match.UTC(), puuid)

	if err != nil {
		log.Println(err)
	}
}

func (leagueDB *LeagueDB) UpdateLastProcessed(lobbyID string) {
	query := `
			UPDATE lobbies
			SET last_processed = ?
			WHERE id = ?`

	_, err := leagueDB.DB.Exec(query, time.Now().UTC(), lobbyID)

	if err != nil {
		log.Println(err)
	}
}

func (leagueDB *LeagueDB) GetPayersMatchesCount(puuid string) int {
	var count int

	query := `
			SELECT COUNT(*) FROM matches
			WHERE player_puuid = ?`

	row := leagueDB.DB.QueryRow(query, puuid)

	row.Scan(&count)

	return count
}
