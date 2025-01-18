package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

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

	fmt.Printf("%v\n", configMap)

	return configMap
}

var Config = getConfig()
