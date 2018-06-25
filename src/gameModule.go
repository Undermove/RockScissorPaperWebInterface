package main

type GameModule struct {
	loosersMap map[string]string
}

func NewGameModule() *GameModule {

	var loosersMap = map[string]string{
		"Rock":     "Scissors",
		"Paper":    "Rock",
		"Scissors": "Paper",
	}

	return &GameModule{
		loosersMap: loosersMap,
	}
}

func (g *GameModule) Turn(playerOne *Player, playerTwo *Player) string {
	var result string

	if playerOne.Choise == playerTwo.Choise {
		result = "DRAW"
	} else if g.loosersMap[playerOne.Choise] == playerTwo.Choise {
		playerOne.Score++
		result = playerOne.Name
	} else {
		playerTwo.Score++
		result = playerTwo.Name
	}

	playerOne.Choise = ""
	playerTwo.Choise = ""

	return result
}
