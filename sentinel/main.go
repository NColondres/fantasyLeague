package main

import (
	"fmt"
	leaguedb "sentinel/internal/leaguedb"
	riotapi "sentinel/internal/riotapi"
	"time"
)

func leagueDBInit() leaguedb.LeagueDB {
	db := leaguedb.LeagueDB{}

	db.SetConfig()
	db.ConnectToDB()
	db.GetMultipliers()

	return db
}

func sentinel() {
	db := leagueDBInit()
	defer db.DB.Close()

	lobbies := db.GetLobbiesByLastProcessed()

	fmt.Println("\nSTARTING NEW SENTINEL\n")

	for _, lobby := range lobbies {

		fmt.Println(lobby.Id, "	", lobby.Last_processed)

		players := db.GetPlayersInLobby(lobby.Id)
		for _, player := range players {
			fmt.Println(player.Name, player.Puuid, player.Last_match.Unix())
			fmt.Println("\nGetting Matches")

			matches := riotapi.GetMatches(player.Puuid, player.Last_match)
			fmt.Printf("\n%+v\n\n", matches)

		}
	}

	fmt.Println("\nSENTINEL COMPLETE")
}

func main() {

	// Run the main function immediately before running periodically at set intervals.
	sentinel()
	for range time.Tick(30 * time.Second) {
		sentinel()
	}

}
