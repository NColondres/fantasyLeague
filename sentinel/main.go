package main

import (
	"fmt"
	"log"
	"reflect"
	leaguedb "sentinel/internal/leaguedb"
	riotapi "sentinel/internal/riotapi"
	"sentinel/internal/score"
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

	for _, lobby := range lobbies {

		players := db.GetPlayersInLobby(lobby.Id, true)

		for _, player := range players {

			// Count the amount of matches the player already played.
			// If they already have played the required matches when the lobby was started, they are skipped.
			playerMatchesCount := db.GetPayersMatchesCount(player.Puuid)

			if !player.Completed {

				// Get all matches the player has played since the last_match timestamp
				fmt.Printf("\nGetting Matches for %s\n", player.Name)

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

								// Keep updating the the last_match_timestamp so that when we are at the end of the loop, we have the most recent gameEndTimestamp.
								// We can then update the players last_match column with the most recent match.
								last_match_timestamp = match_info.GameEndTimeStamp.Add(time.Minute * 2)

							}
						} else {

							// Once the count is equal to the amount of matches set in the lobby, Set the player as completed and break out of the loop.
							db.SetPlayerCompleted(player.Puuid)
							log.Printf("%s has completed all %d matches\n", player.Name, lobby.Matches)

							break
						}
					}

					if !last_match_timestamp.IsZero() {

						db.UpdateLastMatch(player.Puuid, last_match_timestamp)

					}

				}
			}
		}

		// Update lobbies' last_processed timestamp
		db.UpdateLastProcessed(lobby.Id)
	}
}

func scoreCalculate() {

	db := leagueDBInit()
	defer db.DB.Close()

	fmt.Printf("Beginning score calculations..\n\n")

	lobbies := db.GetStartedLobbies()

	for _, lobby := range lobbies {

		players := db.GetPlayersInLobby(lobby, false)

		for _, player := range players {

			scoreCalculator := score.InitScoreCalculator(&db, player, lobby)
			scoreCalculator.Calculate(&db)
		}
	}
}

func main() {

	// Run the main function immediately before running periodically at set intervals.
	sentinel()
	for range time.Tick(30 * time.Second) {
		sentinel()
		go scoreCalculate()
	}
}
