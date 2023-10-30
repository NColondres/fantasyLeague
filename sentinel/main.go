package main

import (
	"fmt"
	"log"
	"reflect"
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

		players := db.GetPlayersInLobby(lobby.Id)

		for _, player := range players {

			// Count the amount of matches the player already played.
			// If they already have played the required matches when the lobby was started, they are skipped.
			playerMatchesCount := db.GetPayersMatchesCount(player.Puuid)
			if playerMatchesCount < lobby.Matches {

				fmt.Printf("\nGetting Matches for %s\n", player.Name)

				// Get all matches the player has played since the last_match timestamp
				matches := riotapi.GetMatches(player.Puuid, player.Last_match)

				if len(matches) == 0 {
					log.Printf("No matches found for %s\n", player.Name)
				} else {

					var last_match_timestamp time.Time

					// Keep count of matches being parsed
					matchCount := playerMatchesCount

					// Loop through each match and get it's info
					for _, match := range matches {

						if matchCount < lobby.Matches {

							match_info := riotapi.GetMatchInfo(match, player.Puuid)

							if !reflect.ValueOf(match_info.Champion).IsZero() {

								fmt.Println(match_info)
								//Insert the grabbed match info into the matches table
								db.InsertMatchInfo(match_info)

								matchCount++
								fmt.Println("matchCount:", matchCount)

								// Keep updating the the last_match_timestamp so that when we are at the end of the loop, we have the most recent gameEndTimestamp.
								// We can then update the players last_match column with the most recent match.
								last_match_timestamp = match_info.GameEndTimeStamp.Add(time.Minute * 2)
							}
						}
					}
					db.UpdateLastMatch(player.Puuid, last_match_timestamp)
				}
			}
		}

		// Update lobbies last_processed timestamp
		db.UpdateLastProcessed(lobby.Id)
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
