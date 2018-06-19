package main

var loosersMap = map[string]string{
	"Rock":     "Scissors",
	"Paper":    "Rock",
	"Scissors": "Paper",
}

type GameModule struct {
}

func turn(playername string, playerChoise string) {
	player := players[playername]
	player.SetPlayerChoise(playerChoise)

	room := rooms[player.CurrentRoomName]
	for _, currentPlayer := range room.Players {
		if currentPlayer.PlayerChoise == "" {
			return
		}
	}

	var result string

	if room.Players[0].PlayerChoise == room.Players[1].PlayerChoise {
		result = "DRAW"
	} else if loosersMap[room.Players[0].PlayerChoise] == room.Players[1].PlayerChoise {
		result = room.Players[0].Name + " WINS!!!"
	} else {
		result = room.Players[1].Name + " WINS!!!"
	}

	for _, currentPlayer := range room.Players {
		authConnections[currentPlayer.Name].Send(result)
	}
}
