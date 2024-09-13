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

func GetSummonerByPUUID(puuid string) map[string]any {
	var puuid_response map[string]any

	requestURI := fmt.Sprintf("https://na1.%s/lol/summoner/v4/summoners/by-puuid/%s", config["BASE_URL"], puuid)
	req, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("X-Riot-Token", config["API_KEY"])

	// Get the response while handling limit request
	res := rateLimitRequests(req)

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode == 200 {

		err2 := json.Unmarshal(resBody, &puuid_response)
		if err2 != nil {
			log.Fatal(err)
		} else {
			return puuid_response
		}
	} else {
		log.Println(string(resBody))
	}

	return puuid_response
}

func GetLeagueAccount(gameName string, tagLine string) map[string]any {

	var account_response map[string]any

	requestURI := fmt.Sprintf("https://americas.%s/riot/account/v1/accounts/by-riot-id/%s/%s", config["BASE_URL"], gameName, tagLine)
	req, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("X-Riot-Token", config["API_KEY"])

	// Get the response while handling limit request
	res := rateLimitRequests(req)

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode == 200 {

		err2 := json.Unmarshal(resBody, &account_response)
		if err2 != nil {
			log.Fatal(err)
		} else {

			// Combine account and summoner info
			combined_info := make(map[string]any)

			summoner_response := GetSummonerByPUUID(account_response["puuid"].(string))

			for key, value := range account_response {
				combined_info[key] = value
			}
			for key, value := range summoner_response {
				combined_info[key] = value
			}

			combined_info["name"] = account_response["gameName"].(string) + "#" + account_response["tagLine"].(string)

			log.Println(combined_info)

			return combined_info
		}
	} else {
		log.Println(string(resBody))
	}
	return account_response
}
