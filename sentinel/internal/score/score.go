package score

import (
	"math"
	"sentinel/internal/leaguedb"
)

type ScoreCalculator struct {
	Player                   leaguedb.Player
	Lobby                    string
	MatchesInfo              []leaguedb.MatchInfo
	Lobby_points_multipliers leaguedb.Lobby_points_multipliers
}

func InitScoreCalculator(leagueDB *leaguedb.LeagueDB, player leaguedb.Player, lobby string) ScoreCalculator {

	var scoreCalculator ScoreCalculator
	scoreCalculator.Player = player
	scoreCalculator.Lobby = lobby
	scoreCalculator.MatchesInfo = leagueDB.GetPlayerMatchesInLobby(scoreCalculator.Player.Puuid, scoreCalculator.Lobby)
	scoreCalculator.Lobby_points_multipliers = leagueDB.GetLobbyPoints(scoreCalculator.Lobby)

	return scoreCalculator

}

func (scoreCalculator *ScoreCalculator) Calculate(leagueDB *leaguedb.LeagueDB) {

	currentScore := scoreCalculator.Player.Total_score

	for _, match := range scoreCalculator.MatchesInfo {

		//Calculate score for KDA separate as more logic is required than the others.m
		scoreCalculator.Player.Total_score += calculateKDA(match.Kills, match.Deaths, match.Assists, scoreCalculator.Lobby_points_multipliers.K_D_A)

		//Convert boolean to float64 to multiply win multiplier by 1 (win) or 0 (loss)
		var win int = 0
		if match.Win {
			win = 1
		}
		//Map the objective to the multiplier for calculations
		objMap := map[string][]float64{
			"Barons":       {float64(match.Barons), float64(scoreCalculator.Lobby_points_multipliers.Baron)},
			"Inhibs":       {float64(match.Inhibs), float64(scoreCalculator.Lobby_points_multipliers.Inhib)},
			"Dragons":      {float64(match.Dragons), float64(scoreCalculator.Lobby_points_multipliers.Dragon)},
			"Rifts":        {float64(match.Rifts), float64(scoreCalculator.Lobby_points_multipliers.Rift)},
			"Turrets":      {float64(match.Turrets), float64(scoreCalculator.Lobby_points_multipliers.Turret)},
			"Pentas":       {float64(match.Pentas), float64(scoreCalculator.Lobby_points_multipliers.Penta)},
			"Quadras":      {float64(match.Quadras), float64(scoreCalculator.Lobby_points_multipliers.Quadra)},
			"Triples":      {float64(match.Triples), float64(scoreCalculator.Lobby_points_multipliers.Triple)},
			"Doubles":      {float64(match.Doubles), float64(scoreCalculator.Lobby_points_multipliers.Double)},
			"Vision_Score": {float64(match.Vision_Score), float64(scoreCalculator.Lobby_points_multipliers.Vision)},
			"Creeps":       {float64(match.Creep_Score), float64(scoreCalculator.Lobby_points_multipliers.Creep)},
			"Win":          {float64(win), float64(scoreCalculator.Lobby_points_multipliers.Win)},
		}

		scoreCalculator.Player.Total_score += calculateObjectives(objMap)

		//Update the current match to be calculated
		leagueDB.UpdateMatchCalculated(scoreCalculator.Player.Puuid, scoreCalculator.Lobby, match.Match_ID)
	}

	//Update the player's total_score only if there is a change
	newScore := scoreCalculator.Player.Total_score

	if currentScore != newScore {
		leagueDB.UpdatePlayerTotalScore(scoreCalculator.Player.Puuid, newScore)
	}

}

//Below will be the functions used to calculate each metric

func calculateKDA(kills int, deaths int, assists int, kdaMultiplier int) int {
	score := 0

	if kills+assists > 0 && kills+assists > deaths {
		//Divide by deaths
		if deaths == 0 {
			score = kills + assists*kdaMultiplier
		} else if deaths > 0 {
			score = int(math.Round(float64(kills+assists)/float64(deaths))) * kdaMultiplier
		}
	}

	return score
}

func calculateObjectives(objective map[string][]float64) int {
	score := 0

	for key, value := range objective {
		points := int(value[0] * value[1])
		score += points
	}

	return score
}
