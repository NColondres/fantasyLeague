package leaguedb

import (
	"database/sql"
	"fmt"
	"log"
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
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db.Config["MYSQL_USER"], db.Config["MYSQL_PASSWORD"], db.Config["MYSQL_URL"], db.Config["MYSQL_PORT"], db.Config["MYSQL_DATABASE"])
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
	ORDER BY last_processed;`

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
	} else {

		for _, lobby := range lobbies {
			fmt.Printf("%+v\n", lobby.Id)
		}

	}

	return lobbies
}
