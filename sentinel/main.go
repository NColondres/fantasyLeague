package main

import (
	"log"
	leaguedb "sentinel/internal/leaguedb"
	"time"
)

func leagueDBInit() leaguedb.LeagueDB {
	db := leaguedb.LeagueDB{}

	db.SetConfig()
	db.ConnectToDB()
	db.GetMultipliers()

	return db
}

func main() {

	db := leagueDBInit()
	defer db.DB.Close()

	for range time.Tick(10 * time.Second) {
		lobbies := db.GetLobbiesByLastProcessed()

		log.Println(lobbies)
	}
}
