package riotapi

import (
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
