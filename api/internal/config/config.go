package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Read .env file
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

// Read multipliers for scoring.toml file located in root of project
func getScoring() map[string]any {

	viper.SetConfigFile("scoring.toml")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	scoringMap := make(map[string]any)

	for _, v := range viper.AllKeys() {
		scoringMap[v] = viper.Get(v)
	}

	return scoringMap
}

func init() {

	Config = getConfig()
	Scoring = getScoring()

}

// Exporting variables

var Config map[string]string

var Scoring map[string]any
