package riotapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type MatchInfo struct {
	Match_ID         string
	Puuid            string
	Champion         string
	Position         string
	Kills            int
	Deaths           int
	Assists          int
	Turrets          int
	Inhibs           int
	Dragons          int
	Rifts            int
	Barons           int
	Vision_Score     int
	Creep_Score      int
	Pentas           int
	Quadras          int
	Triples          int
	Doubles          int
	Win              bool
	GameEndTimeStamp time.Time
}

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

// Helper function to manage rate limits of Riot API
func rateLimitRequests(req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if res.StatusCode == 403 {
		log.Fatalln("RIOT API: API KEY IS INVALID.")
	} else if res.StatusCode == 429 {
		resBody, _ := io.ReadAll(res.Body)
		retry_after, _ := strconv.Atoi(res.Header.Get("Retry-After"))
		log.Printf("Retrying after %d seconds: %v", retry_after, string(resBody))
		time.Sleep(time.Duration(retry_after) * time.Second)

		return rateLimitRequests(req)

	}
	return res
}

//End of helper functions

func GetMatches(puuid string, last_match *time.Time) []string {

	requestURI := fmt.Sprintf("https://%s.%s/lol/match/v5/matches/by-puuid/%s/ids?startTime=%d", "americas", config["BASE_URL"], puuid, last_match.Unix())
	req, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Riot-Token", config["API_KEY"])

	// Get the response while handling limit request
	res := rateLimitRequests(req)

	var response []string
	if res.StatusCode == 200 {
		resBody, _ := io.ReadAll(res.Body)

		err := json.Unmarshal(resBody, &response)
		if err != nil {
			log.Println(err)
		}
		sort.Strings(response)
	}
	return response
}

func GetMatchInfo(matchID string, puuid string) MatchInfo {

	matchInfo := MatchInfo{
		Match_ID: matchID,
		Puuid:    puuid,
	}

	requestURI := fmt.Sprintf("https://%s.%s/lol/match/v5/matches/%s", "americas", config["BASE_URL"], matchID)
	req, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Riot-Token", config["API_KEY"])

	// Get the response while handling limit request
	res := rateLimitRequests(req)

	var response map[string]any
	if res.StatusCode == 200 {
		resBody, _ := io.ReadAll(res.Body)

		err := json.Unmarshal(resBody, &response)
		if err != nil {
			log.Println(err)
		}
	}

	// Check to make sure it is a summoner's rift game. Can be ranked or normals
	gameMode := response["info"].(map[string]any)["gameMode"].(string)
	queueID := int(response["info"].(map[string]any)["queueId"].(float64))

	if gameMode == "CLASSIC" && queueID >= 400 && queueID <= 440 {

		participants := response["info"].(map[string]any)["participants"].([]any)

		for _, participant := range participants {

			// Find the enrolled player and start grabbing desired info
			participant_info := participant.(map[string]any)
			if participant_info["puuid"] == puuid {

				// Populate struct with values from the json response
				matchInfo.Champion = participant_info["championName"].(string)
				matchInfo.Position = participant_info["teamPosition"].(string)
				matchInfo.Kills = int(participant_info["kills"].(float64))
				matchInfo.Deaths = int(participant_info["deaths"].(float64))
				matchInfo.Assists = int(participant_info["assists"].(float64))
				matchInfo.Turrets = int(participant_info["turretTakedowns"].(float64))
				matchInfo.Inhibs = int(participant_info["inhibitorTakedowns"].(float64))
				matchInfo.Dragons = int(participant_info["dragonKills"].(float64))
				matchInfo.Rifts = int(participant_info["challenges"].(map[string]any)["riftHeraldTakedowns"].(float64))
				matchInfo.Barons = int(participant_info["baronKills"].(float64))
				matchInfo.Vision_Score = int(participant_info["visionScore"].(float64))
				matchInfo.Creep_Score = int(participant_info["totalMinionsKilled"].(float64)) + int(participant_info["neutralMinionsKilled"].(float64))
				matchInfo.Pentas = int(participant_info["pentaKills"].(float64))
				matchInfo.Quadras = int(participant_info["quadraKills"].(float64))
				matchInfo.Triples = int(participant_info["tripleKills"].(float64))
				matchInfo.Doubles = int(participant_info["doubleKills"].(float64))
				matchInfo.Win = participant_info["win"].(bool)
				matchInfo.GameEndTimeStamp = time.Unix(0, int64(response["info"].(map[string]any)["gameEndTimestamp"].(float64))*int64(time.Millisecond)).UTC()
			}
		}
	}

	return matchInfo

}
