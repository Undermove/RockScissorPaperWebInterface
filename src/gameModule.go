package main

var loosersMap = map[string]string{
	"Rock":     "Scissors",
	"Paper":    "Rock",
	"Scissors": "Paper",
}

type GameModule struct {
}

func turn(playerOneChoise string, playerTwoChoice string) string {
	var result string

	if playerOneChoise == playerTwoChoice {
		result = "DRAW"
	} else if loosersMap[playerOneChoise] == playerTwoChoice {
		result = "PlayerOne WINS!!!"
	} else {
		result = "PlayerTwo WINS!!!"
	}

	return result
}
