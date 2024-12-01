package main

import (
	leaguedb "api/internal/leaguedb"
	"api/internal/riotapi"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type player struct {
	GameName string `json:"gameName" binding:"required"`
	Tag_Line string `json:"tagLine" binding:"required"`
	Lobby_Id string `json:"lobby_id"`
}

type cookies struct {
	lobby_id string
	puuid    string
}

type start_lobby struct {
	Lobby   leaguedb.Lobby    `json:"lobby"`
	Players []leaguedb.Player `json:"players"`
}

func getCookies(c *gin.Context) cookies {
	lobby_id, _ := c.Cookie("lobby_id")
	puuid_cookie, _ := c.Cookie("puuid")
	return cookies{
		lobby_id: lobby_id,
		puuid:    puuid_cookie,
	}
}

// func serverHeader(c *gin.Context) {

// 	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
// 	c.Header("Access-Control-Allow-Credentials", "true")
// 	c.Header("Access-Control-Allow-Methods", "GET,HEAD,POST")
// }

func main() {

	// Gin API Server
	router := gin.Default()

	// Set the api headers needed for the web app to use the api appropriately
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4173"},
		AllowMethods:     []string{"GET", "HEAD", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Healthy": 200})
	})

	router.GET("/player", func(c *gin.Context) {
		c.JSON(http.StatusOK, leaguedb.GetPlayers())
	})

	router.DELETE("/player/:puuid", func(c *gin.Context) {

		// url param
		puuid := c.Param("puuid")

		// Get the cookies of the requestor
		cookies := getCookies(c)

		// Check to make sure cookies puuid and lobby_id are set
		if cookies.lobby_id == "" || cookies.puuid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Must have both coockies 'lobby_id' and 'puuid set'"})
		} else {

			lobby, lobbyErr := leaguedb.GetLobby(cookies.lobby_id)

			player, playerErr := leaguedb.GetPlayerInLobby(cookies.lobby_id, puuid)

			// Validation stuff
			switch {

			// Check if lobby exists in db
			case lobbyErr != nil:
				c.JSON(http.StatusBadRequest, gin.H{"error": "lobby not found"})

			// Check if cookie.puuid is the creator of lobby if not, they cannot removed someone from the lobby
			case cookies.puuid != lobby.Creator_puuid:
				c.JSON(http.StatusBadRequest, gin.H{"denied": fmt.Sprintf("%s is not the creator of lobby: %s", cookies.puuid, lobby.Id)})

			// Check if player exists in lobby
			case playerErr != nil:
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s not found in lobby: %s", puuid, cookies.lobby_id)})

			// Check if player has completed all required games. Players who have completed all their games cannot be deleted from the lobby.
			case player.Completed:
				c.JSON(http.StatusBadRequest, gin.H{"denied": fmt.Sprintf("%s has completed all their games. Cannot delete", player.Name)})

			default:
				c.JSON(http.StatusOK, leaguedb.DeletePlayerFromLobby(player.Puuid, cookies.lobby_id))
			}

		}

	})

	router.POST("/start", func(c *gin.Context) {
		// Confirm and grab cookies
		cookies := getCookies(c)
		if cookies.lobby_id == "" || cookies.puuid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Must have both cookies 'lobby_id' and 'puuid set'"})
		} else {
			// Get Lobby record, if cookie.puuid is equal to lobby record creator_puuid and the , that means they are allowed to start the tournament for everyone enrolled.
			lobby, _ := leaguedb.GetLobby(cookies.lobby_id)

			if lobby.Started {
				c.JSON(http.StatusBadRequest, gin.H{"error": "lobby already started"})

			} else if cookies.puuid == lobby.Creator_puuid {

				// Get body from request
				type body struct {
					Rules                    map[string]any
					Lobby_points_multipliers leaguedb.Lobby_points_multipliers
				}

				var response body
				c.ShouldBindJSON(&response)

				// Update lobby struct with values given in request body
				lobby.Lobby_points_multipliers = response.Lobby_points_multipliers

				// Set lobby points
				err := leaguedb.SetLobbyPoints(cookies.lobby_id, &lobby.Lobby_points_multipliers)

				if err != nil {

					log.Println(err)

				} else {

					//Start the lobby
					lobby_start, start_lobby_err := leaguedb.StartLobby(cookies.lobby_id, response.Rules)

					if start_lobby_err != nil {

						c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(start_lobby_err)})

					} else {

						lobby, _ = leaguedb.GetLobby(cookies.lobby_id)
						lobby.Lobby_points_multipliers = leaguedb.GetLobbyPoints(cookies.lobby_id)
						start_lobby := start_lobby{
							Lobby:   lobby,
							Players: lobby_start,
						}
						c.JSON(http.StatusOK, start_lobby)
					}
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"denied": "Must be creator of Lobby to start it."})
			}
		}

	})

	// Handle enrolling players. Requires at minumum the 'summoner' name of the user enrolling in json.
	router.POST("/enroll", func(c *gin.Context) {

		c.SetSameSite(http.SameSiteStrictMode)

		lobby_id_cookie, lobby_id_cookie_err := c.Cookie("lobby_id")
		_, puuid_cookie_err := c.Cookie("puuid")

		newPlayer := player{}

		if lobby_id_cookie_err == nil && puuid_cookie_err == nil {
			c.JSON(http.StatusForbidden, gin.H{"denied": "cookies show player is already enrolled"})

		} else if err := c.BindJSON(&newPlayer); err != nil {

			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})

		} else if lobby, _ := leaguedb.GetLobby(lobby_id_cookie); lobby.Started {
			c.JSON(http.StatusForbidden, gin.H{"denied": "lobby has already started"})
		} else {
			// Get the league account information from Riot API
			log.Println(newPlayer)
			summoner_account := riotapi.GetLeagueAccount(newPlayer.GameName, newPlayer.Tag_Line)
			// Check if summoner map is empty.
			if len(summoner_account) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s could not be found", newPlayer.GameName)})
			} else {

				// Append lobby_id if cookie exits
				if lobby_id_cookie_err == nil {
					summoner_account["lobby_id"] = lobby_id_cookie
				}
				// Add user to database
				lobby_id, leagueErr := leaguedb.AddSummoner(summoner_account)
				if leagueErr != nil {
					log.Println(leagueErr)
					c.JSON(http.StatusBadRequest, gin.H{"error": "failed to add player"})

				} else {
					// Set cookie if lobby_id has been generated
					if lobby_id_cookie_err != nil && lobby_id != "" {
						c.SetCookie("lobby_id", lobby_id, 3600*24*7, "/", "localhost", false, false)
					}

					// Set puuid cookie to track which user is making requests
					if puuid_cookie_err != nil {
						c.SetCookie("puuid", summoner_account["puuid"].(string), 3600*24*7, "/", "localhost", false, false)
					}
					c.JSON(http.StatusCreated, gin.H{"success": summoner_account["name"]})
				}

			}

		}

	})

	//Get status of a lobby
	router.GET("/lobby", func(c *gin.Context) {

		lobby_id := getCookies(c).lobby_id

		var players []leaguedb.Player

		if lobby_id != "" {
			lobby, lobbyErr := leaguedb.GetLobby(lobby_id)

			if lobbyErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "lobby not found"})
			} else {

				players = leaguedb.GetPlayersInLobby(lobby_id)
				lobby.Lobby_points_multipliers = leaguedb.GetLobbyPoints(lobby_id)

				// Add matches info for each player if the lobby.Started == true when this endpoint is requested
				if lobby.Started {

					for index := range players {

						players[index].Matches = leaguedb.GetPlayerMatchesInLobby(players[index].Puuid, lobby.Id)

					}
				}

				lobby_players := map[string]any{
					"lobby":   lobby,
					"players": players,
				}

				c.JSON(http.StatusOK, lobby_players)

			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no lobby_id cookie set"})
		}
	})

	router.GET("/player/matches/:puuid", func(c *gin.Context) {

		//url param
		puuid := c.Param("puuid")
		cookies := getCookies(c)

		if cookies.lobby_id == "" {

			c.JSON(http.StatusBadRequest, gin.H{"Error": "missing lobby_id cookie"})

		} else {

			matches := leaguedb.GetPlayerMatchesInLobby(puuid, cookies.lobby_id)

			c.JSON(http.StatusOK, matches)
		}

	})

	router.Run()

}
