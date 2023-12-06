package leagueDB

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// Read the .env file in the root directory.
var config = getConfig()

func getConfig() map[string]string {
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
	return configMap
}

// --- Bunch of helper functions ---

// Read multipliers for scoring.toml file located in root of project
func getMultipliers() map[string]any {
	viper.SetConfigFile("scoring.toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	multipliersMap := make(map[string]any)
	for _, v := range viper.AllKeys() {
		multipliersMap[v] = viper.Get(v)
	}
	return multipliersMap
}
func connectToDB() *sql.DB {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config["MYSQL_USER"], config["MYSQL_PASSWORD"], config["MYSQL_URL"], config["MYSQL_PORT"], config["MYSQL_DATABASE"])
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("Error while connecting to MySQL DB %s", err)
	}
	return db
}

func randStr() string {
	// rand.Seed(time.Now().Unix())
	alphanumeric := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, 17)
	for i := range b {
		// randomly select 1 character from given charset
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

// --- End of help functions ---

type PlayerAll struct {
	Puuid           string      `json:"puuid"`
	Id              string      `json:"id"`
	Account_id      string      `json:"account_id"`
	Name            string      `json:"name"`
	Profile_icon_id int         `json:"profile_icon_id"`
	Level           int         `json:"level"`
	Revision_date   int         `json:"revision_date"`
	Region          string      `json:"region"`
	Lobby_Id        *string     `json:"lobby_id"`
	Last_match      *time.Time  `json:"last_match"`
	Total_score     int         `json:"total_score"`
	Completed       *bool       `json:"completed"`
	Matches         []MatchInfo `json:"matches,omitempty"`
}

type Player struct {
	Puuid           string      `json:"puuid"`
	Name            string      `json:"name"`
	Level           int         `json:"level"`
	Region          string      `json:"region"`
	Last_match      *time.Time  `json:"last_match"`
	Total_score     int         `json:"total_score"`
	Completed       bool        `json:"completed"`
	Matches         []MatchInfo `json:"matches,omitempty"`
	Profile_icon_id int         `json:"profile_icon_id"`
}

type Lobby struct {
	Id                       string                   `json:"id"`
	Creation_time            time.Time                `json:"creation_time"`
	Creator_puuid            string                   `json:"creator_puuid"`
	Started                  bool                     `json:"started"`
	Start_time               *time.Time               `json:"started_time,omitempty"`
	Matches                  int                      `json:"matches"`
	Last_processed           time.Time                `json:"last_processed"`
	Lobby_points_multipliers Lobby_points_multipliers `json:"lobby_points_multipliers,omitempty"`
}

type Lobbies_Effected struct {
	Player_puuid       string     `json:"player_puuid"`
	Player_name        string     `json:"player_name"`
	Player_last_match  *time.Time `json:"player_last_match"`
	Lobby_id           string     `json:"lobby_id"`
	Lobby_started      bool       `json:"lobby_started"`
	Lobby_started_time *time.Time `json:"lobby_started_time"`
}

type Lobby_points_multipliers struct {
	K_d_a  int     `json:"k_d_a"`
	Baron  int     `json:"baron"`
	Inhib  int     `json:"inhib"`
	Dragon int     `json:"dragon"`
	Turret int     `json:"turret"`
	Penta  int     `json:"penta"`
	Quadra int     `json:"quadra"`
	Triple int     `json:"triple"`
	Double int     `json:"double"`
	Rift   int     `json:"rift"`
	Vision int     `json:"vision"`
	Win    int     `json:"win"`
	Creep  float32 `json:"creep"`
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
}

func GetPlayers() []PlayerAll {
	db := connectToDB()
	defer db.Close()
	results, err := db.Query("SELECT * FROM players LIMIT 100;")
	if err != nil {
		log.Fatal(err)
	}

	var players []PlayerAll

	for results.Next() {

		var player PlayerAll

		err := results.Scan(&player.Puuid, &player.Id, &player.Account_id, &player.Profile_icon_id,
			&player.Revision_date, &player.Level, &player.Name, &player.Region, &player.Lobby_Id, &player.Last_match, &player.Total_score, &player.Completed)
		if err != nil {
			log.Fatal(err)
		}

		players = append(players, player)
	}

	results.Close()

	return players
}

func GetPlayersInLobby(lobby_id string) []Player {
	db := connectToDB()
	defer db.Close()

	query := `
			SELECT players.puuid, players.name, players.level, players.region, players.last_match, players.total_score, players.completed, players.profile_icon_id
			FROM players
			JOIN lobbies ON players.lobby_id = lobbies.id
			WHERE lobbies.id = ?
			ORDER BY players.total_score DESC;`

	results, err := db.Query(query, lobby_id)
	if err != nil {
		log.Fatal(err)
	}

	var players []Player

	for results.Next() {

		var player Player

		err := results.Scan(&player.Puuid, &player.Name, &player.Level, &player.Region, &player.Last_match, &player.Total_score, &player.Completed, &player.Profile_icon_id)
		if err != nil {
			log.Fatal(err)
		}

		players = append(players, player)
	}

	results.Close()
	return players

}

func GetPlayerInLobby(lobby_id string, puuid string) (Player, error) {
	db := connectToDB()
	defer db.Close()
	query := `
			SELECT players.puuid, players.name, players.level, players.region, players.last_match, players.completed
			FROM players
			JOIN lobbies ON players.lobby_id = lobbies.id
			WHERE (lobbies.id = ? AND players.puuid = ?);`

	var player Player

	results := db.QueryRow(query, lobby_id, puuid)

	err := results.Scan(&player.Puuid, &player.Name, &player.Level, &player.Region, &player.Last_match, &player.Completed)

	if err != nil {
		return player, err
	}
	return player, nil
}

func GetPlayerMatchesInLobby(puuid string, lobbyID string) []MatchInfo {
	// Connect to DB
	db := connectToDB()
	defer db.Close()

	var matches []MatchInfo

	// Getting match info from matches table but only from the lobby requested.
	query := `
	SELECT match_id, player_puuid, champion, position, kills, deaths, assists, turrets, inhibs, dragons, rifts, barons, vision_score, creep_score, pentas, quadras, triples, doubles, win, game_end_timestamp
	FROM matches
	JOIN players ON players.puuid = matches.player_puuid
	WHERE player_puuid = ? AND players.lobby_id = ?;`

	rows, err := db.Query(query, puuid, lobbyID)

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {

		var matchInfo MatchInfo

		err := rows.Scan(&matchInfo.Match_ID, &matchInfo.Puuid, &matchInfo.Champion,
			&matchInfo.Position, &matchInfo.Kills, &matchInfo.Deaths, &matchInfo.Assists, &matchInfo.Turrets,
			&matchInfo.Inhibs, &matchInfo.Dragons, &matchInfo.Rifts, &matchInfo.Barons, &matchInfo.Vision_Score,
			&matchInfo.Creep_Score, &matchInfo.Pentas, &matchInfo.Quadras, &matchInfo.Triples, &matchInfo.Doubles, &matchInfo.Win, &matchInfo.Game_End_Timestamp)

		if err != nil {
			log.Fatal(err)
		}

		matches = append(matches, matchInfo)
	}

	return matches

}

func AddSummoner(player map[string]any) (string, error) {
	// Connect to DB
	db := connectToDB()
	defer db.Close()

	// Player does not have a lobby_id which means they will be making a new one and joining it as the creator.

	// Create a new non-existing lobby id
	if _, exist := player["lobby_id"]; !exist {

		lobby_id := struct {
			id string
		}{
			randStr(),
		}
		lobby_result := db.QueryRow("SELECT id FROM lobbies WHERE id = ?", lobby_id.id)

		// While Loop to generate a new lobby until a unique one is available.
		for lobby_result.Scan() == nil {
			log.Printf("Lobby: %s already exists, generating new lobby\n", lobby_id.id)
			lobby_id.id = randStr()
			lobby_result = db.QueryRow("SELECT id FROM lobbies WHERE id = ?", lobby_id.id)

		}
		// Add lobby id to player map
		player["lobby_id"] = lobby_id.id

		// Create a new lobby with confirmed available lobby id
		_, err := db.Exec("INSERT INTO lobbies (id, creator_puuid) VALUES (?,?)", player["lobby_id"], player["puuid"])
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Printf("%s is being added to lobby: %s\n", player["name"], player["lobby_id"])
	}

	// Now add player to players table
	query, err := db.Prepare("INSERT INTO players (puuid, id, account_id, profile_icon_id, revision_date, level, name, region, lobby_id) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return "", err
	}

	_, execErr := query.Exec(player["puuid"], player["id"], player["accountId"],
		player["profileIconId"], player["revisionDate"],
		player["summonerLevel"], player["name"], player["region"], player["lobby_id"])
	if execErr != nil {
		return "", execErr
	}
	query.Close()

	return player["lobby_id"].(string), nil
}

func GetLobby(lobby_id string) (Lobby, error) {
	db := connectToDB()
	defer db.Close()

	var lobby Lobby

	results := db.QueryRow("SELECT id, creation_time, creator_puuid, started, started_time, matches, last_processed FROM lobbies WHERE id = ? LIMIT 1;", lobby_id)

	if results.Err() != nil {
		log.Println(results.Err())
		return lobby, results.Err()
	}

	results.Scan(&lobby.Id, &lobby.Creation_time, &lobby.Creator_puuid, &lobby.Started, &lobby.Start_time, &lobby.Matches, &lobby.Last_processed)

	if reflect.ValueOf(lobby).IsZero() {
		return lobby, errors.New("No lobby found")
	}

	return lobby, nil

}

func GetLobbyPoints(lobby_id string) Lobby_points_multipliers {

	db := connectToDB()
	defer db.Close()

	query := `SELECT k_d_a, baron, inhib, dragon, turret, penta, quadra, triple, lobby_points_multipliers.double, rift, vision, win, creep
	FROM lobby_points_multipliers
	WHERE lobby_id = ?;`

	var multipliers Lobby_points_multipliers

	result := db.QueryRow(query, lobby_id)

	if result.Err() == nil {
		result.Scan(&multipliers.K_d_a, &multipliers.Baron, &multipliers.Inhib, &multipliers.Dragon, &multipliers.Turret, &multipliers.Penta,
			&multipliers.Quadra, &multipliers.Triple, &multipliers.Double, &multipliers.Rift, &multipliers.Vision, &multipliers.Win, &multipliers.Creep)
	} else {
		log.Println(result.Err())
	}

	return multipliers

}

func StartLobby(lobby_id string, rules map[string]any) ([]Player, error) {

	db := connectToDB()
	defer db.Close()
	var players []Player
	var err error

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	// Get the amount of enrolled players in lobby
	var count int
	tx.QueryRow("SELECT COUNT(*) FROM players WHERE lobby_id = ?", lobby_id).Scan(&count)

	// Logic to use either default matches or matches passed in through function
	matches, exist := rules["matches"]
	if exist == false || reflect.ValueOf(rules["matches"].(float64)).IsZero() == true {

		matches, _ = strconv.Atoi(config["DEFAULT_MATCHES"])
	} else {
		matches = int(rules["matches"].(float64))
	}

	// Only start lobby if the count is >= MINIMUM_PLAYERS ENV variable
	minimum_players, _ := strconv.Atoi(config["MINIMUM_PLAYERS"])

	if count >= minimum_players {

		// Generate current UTC time. Adding 10 minutes to prevent lobbies from being started during a match.
		current_time := time.Now().UTC().Add(time.Duration(10) * time.Minute)

		// Update lobby's started boolen to true and started_time to current UTC time
		_, lobbies_err := tx.Exec("UPDATE lobbies SET started = ?, started_time = ?, matches = ? WHERE id = ?", true, current_time, matches, lobby_id)
		if lobbies_err != nil {
			log.Fatal(lobbies_err)
		}

		// Update players' in lobby last_match column to current UTC time to begin tracking of matches after that point (Might want to add 10 minutes to prevent being cheeky)
		_, players_err := tx.Exec("UPDATE players SET last_match = ? WHERE lobby_id = ?", current_time, lobby_id)
		if players_err != nil {
			log.Fatal(players_err)
		}

		// Commit changes to DB
		commit_err := tx.Commit()
		if commit_err != nil {
			log.Fatal(commit_err)
		}

		// Query db to return what has been changed by this function
		query := `
				SELECT players.puuid, players.name, players.level, players.region, players.last_match, players.completed
				FROM players
				JOIN lobbies ON players.lobby_id = lobbies.id
				WHERE lobbies.id = ?;`
		results, query_err := db.Query(query, lobby_id)
		if query_err != nil {
			log.Println(query_err)
		} else {

			for results.Next() {
				var player Player
				results.Scan(&player.Puuid, &player.Name, &player.Level, &player.Region, &player.Last_match, &player.Completed)
				players = append(players, player)
			}
			return players, nil
		}

	} else {
		err = fmt.Errorf("Not enough players enrolled to start. Minimum requirement %d", minimum_players)
		return players, err
	}

	return players, err
}

func SetLobbyPoints(lobby_id string, multipliers *Lobby_points_multipliers) error {

	scoringConfig := getMultipliers()
	values := reflect.Indirect(reflect.ValueOf(multipliers))
	types := values.Type()

	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).IsZero() {
			switch struct_type := types.Field(i).Type.Name(); struct_type {
			case "int":
				values.Field(i).SetInt(scoringConfig[strings.ToLower(types.Field(i).Name)].(int64))
			case "float32":
				values.Field(i).SetFloat(scoringConfig[strings.ToLower(types.Field(i).Name)].(float64))
			}
		}
	}
	query := `
			REPLACE INTO lobby_points_multipliers (lobby_id, k_d_a, baron, inhib, dragon, turret, penta, quadra, triple, lobby_points_multipliers.double, rift, vision, win, creep) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	db := connectToDB()
	defer db.Close()

	_, err := db.Exec(
		query, lobby_id, multipliers.K_d_a, multipliers.Baron, multipliers.Inhib,
		multipliers.Dragon, multipliers.Turret, multipliers.Penta, multipliers.Quadra,
		multipliers.Triple, multipliers.Double, multipliers.Rift, multipliers.Vision,
		multipliers.Win, multipliers.Creep)

	if err != nil {
		return err
	}
	return nil
}

func DeletePlayerFromLobby(puuidToDelete string, lobby_id string) map[string]any {

	db := connectToDB()
	defer db.Close()

	var response map[string]any

	deletePlayer, _ := GetPlayerInLobby(lobby_id, puuidToDelete)

	query := `
			DELETE FROM players
			WHERE (puuid = ? AND lobby_id = ?)`

	_, err := db.Exec(query, puuidToDelete, lobby_id)

	if err != nil {
		log.Println(err)
	} else {
		// Get the new list of players in lobby
		lobbyPlayers := GetPlayersInLobby(lobby_id)

		response = map[string]any{
			"deleted": deletePlayer,
			"players": lobbyPlayers,
		}
	}

	return response

}
