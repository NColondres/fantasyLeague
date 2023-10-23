package riotapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

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
	test_api(configMap)
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
		log.Printf("Retrying after: %d seconds\n%v", retry_after, string(resBody))
		time.Sleep(time.Duration(retry_after) * time.Second)

		return rateLimitRequests(req)

	}
	return res
}

// Test to confirm api key is working. If it fails, API service will exit with error

func GetLeagueAccount(summonerName string, summonerRegion string) map[string]any {

	requestURI := fmt.Sprintf("https://%s.%s/lol/summoner/v4/summoners/by-name/%s", summonerRegion, config["BASE_URL"], summonerName)
	req, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("X-Riot-Token", config["API_KEY"])

	// Get the response while handling limit request
	res := rateLimitRequests(req)

	var response map[string]any
	if res.StatusCode == 200 {
		resBody, _ := io.ReadAll(res.Body)

		err2 := json.Unmarshal(resBody, &response)
		if err2 != nil {
			log.Fatal(err)
		} else {
			response["region"] = summonerRegion
			return response
		}

	} else {
		resBody, _ := io.ReadAll(res.Body)
		log.Println(string(resBody))
	}
	return response
}
