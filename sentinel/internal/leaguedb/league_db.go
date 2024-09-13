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
	Region      *string
	Last_match  *time.Time
	Total_score int
	Completed   bool
}

type MatchInfo struct {
	Match_ID           string    `json:"match_id"`
	Puuid              string    `json:"puuid"`
	Champion           string    `json:"champion"`
	Position           string    `json:"position"`
	Kills              int       `json:"kills"`
	Deaths             int       `json:"deaths"`
	Assists            int       `json:"assists"`
	Turrets            int       `json:"turrets"`
	Inhibs             int       `json:"inhibs"`
	Dragons            int       `json:"dragons"`
	Rifts              int       `json:"rifts"`
	Barons             int       `json:"barons"`
	Vision_Score       int       `json:"vision_score"`
	Creep_Score        int       `json:"creep_score"`
	Pentas             int       `json:"pentas"`
	Quadras            int       `json:"quadras"`
	Triples            int       `json:"triples"`
	Doubles            int       `json:"doubles"`
	Win                bool      `json:"win"`
	Game_End_Timestamp time.Time `json:"game_end_timestamp"`
	Calculated         bool      `json:"-"`
}
type Lobby_points_multipliers struct {
	Lobby_ID string
	K_D_A    int
	Turret   int
	Inhib    int
	Dragon   int
	Rift     int
	Baron    int
	Vision   int
	Creep    float32
	Penta    int
	Quadra   int
	Triple   int
	Double   int
	Win      int
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

	pingErr := db.DB.Ping()

	for pingErr != nil {

		log.Printf("Waiting for DB at %s to respond to ping\n", db.Config["MYSQL_URL"])
		time.Sleep(time.Second * 10)
		pingErr = db.DB.Ping()

	}
}

func (db *LeagueDB) GetLobbiesByLastProcessed() []Lobby {

	query := `
	SELECT id, creation_time, creator_puuid, started, started_time, matches, last_processed
	FROM lobbies
	WHERE started = TRUE
	ORDER BY last_processed
	LIMIT 5;`

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

func (leagueDB *LeagueDB) GetPlayersInLobby(lobbyID string, notCompletedOnly bool) []Player {

	players := []Player{}

	query := `
		SELECT puuid, name, region, last_match, total_score, completed
		FROM players
		WHERE lobby_id = ?`

	if notCompletedOnly {

		query = query + " AND completed = FALSE"
	}

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
								dragons, rifts, barons, vision_score, creep_score, pentas, quadras, triples, doubles, win, game_end_timestamp)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err := leagueDB.DB.Exec(
		query, matchInfo.Match_ID, matchInfo.Puuid, matchInfo.Champion, matchInfo.Position, matchInfo.Kills, matchInfo.Deaths,
		matchInfo.Assists, matchInfo.Turrets, matchInfo.Inhibs, matchInfo.Dragons, matchInfo.Rifts, matchInfo.Barons, matchInfo.Vision_Score,
		matchInfo.Creep_Score, matchInfo.Pentas, matchInfo.Quadras, matchInfo.Triples, matchInfo.Doubles, matchInfo.Win, matchInfo.GameEndTimeStamp)

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

func (leagueDB *LeagueDB) SetPlayerCompleted(puuid string) {
	query := `
			UPDATE players
			SET completed = TRUE
			WHERE puuid = ?`

	_, err := leagueDB.DB.Exec(query, puuid)

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

	if row.Err() != nil {
		log.Println(row.Err())
	}

	row.Scan(&count)

	return count
}

func (leagueDB *LeagueDB) GetStartedLobbies() []string {

	var lobbies []string

	query := `
	SELECT id FROM lobbies
	WHERE started = TRUE
	ORDER BY last_processed`

	rows, err := leagueDB.DB.Query(query)

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var lobby string

		rows.Scan(&lobby)

		lobbies = append(lobbies, lobby)
	}

	return lobbies
}

func (leagueDB *LeagueDB) GetLobbyPoints(lobby_id string) Lobby_points_multipliers {

	query := `SELECT k_d_a, baron, inhib, dragon, turret, penta, quadra, triple, lobby_points_multipliers.double, rift, vision, win, creep
	FROM lobby_points_multipliers
	WHERE lobby_id = ?;`

	var multipliers Lobby_points_multipliers

	result := leagueDB.DB.QueryRow(query, lobby_id)

	if result.Err() == nil {
		result.Scan(&multipliers.K_D_A, &multipliers.Baron, &multipliers.Inhib, &multipliers.Dragon, &multipliers.Turret, &multipliers.Penta,
			&multipliers.Quadra, &multipliers.Triple, &multipliers.Double, &multipliers.Rift, &multipliers.Vision, &multipliers.Win, &multipliers.Creep)
	} else {
		log.Println(result.Err())
	}

	return multipliers

}

func (leagueDB *LeagueDB) GetPlayerMatchesInLobby(puuid string, lobbyID string) []MatchInfo {

	var matches []MatchInfo

	// Getting match info from matches table but only from the lobby requested.
	query := `
	SELECT match_id, player_puuid, champion, position, kills, deaths, assists, turrets, inhibs, dragons, rifts, barons, vision_score, creep_score, pentas, quadras, triples, doubles, win, game_end_timestamp, calculated
	FROM matches
	JOIN players ON players.puuid = matches.player_puuid
	WHERE player_puuid = ? AND players.lobby_id = ? AND calculated = FALSE;`

	rows, err := leagueDB.DB.Query(query, puuid, lobbyID)

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {

		var matchInfo MatchInfo

		err := rows.Scan(&matchInfo.Match_ID, &matchInfo.Puuid, &matchInfo.Champion,
			&matchInfo.Position, &matchInfo.Kills, &matchInfo.Deaths, &matchInfo.Assists, &matchInfo.Turrets,
			&matchInfo.Inhibs, &matchInfo.Dragons, &matchInfo.Rifts, &matchInfo.Barons, &matchInfo.Vision_Score,
			&matchInfo.Creep_Score, &matchInfo.Pentas, &matchInfo.Quadras, &matchInfo.Triples, &matchInfo.Doubles, &matchInfo.Win, &matchInfo.Game_End_Timestamp, &matchInfo.Calculated)

		if err != nil {
			log.Fatal(err)
		}

		matches = append(matches, matchInfo)
	}

	return matches

}

func (leagueDB *LeagueDB) UpdateMatchCalculated(puuid string, lobbyID string, matchID string) {

	query := `
			UPDATE matches INNER JOIN players ON (matches.player_puuid = players.puuid)
			SET matches.calculated = TRUE
			WHERE matches.player_puuid = ? AND players.lobby_id = ? AND matches.match_id = ?`

	_, err := leagueDB.DB.Exec(query, puuid, lobbyID, matchID)

	if err != nil {
		log.Println(err)
	}
}

func (leagueDB *LeagueDB) UpdatePlayerTotalScore(puuid string, score int) {

	query := `
	UPDATE players
	SET total_score = ?
	WHERE puuid = ?`

	_, err := leagueDB.DB.Exec(query, score, puuid)

	if err != nil {
		log.Println(err)
	}
}
