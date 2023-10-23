package riotapi

import (
	"fmt"
	"log"
	"net/http"
)

// Quick function to make sure the RIOT Api is working. App will exit with error if unable to get 200 response status
func test_api(config map[string]string) {

	requestURI := fmt.Sprintf("https://na1.%s/lol/status/v4/platform-data", config["BASE_URL"])
	req, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("X-Riot-Token", config["API_KEY"])
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode == 403 {
		log.Fatalln("RIOT API: API KEY IS INVALID.")
	} else if res.StatusCode != 200 {
		log.Fatalln("Riot Api: Something is wrong with Riot's API service")
	}

}
