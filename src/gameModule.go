package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

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

func (rm *RoomsManager) SendResponse(result string, w *websocket.Conn) {
	var response CreateRoomResponse

	if isSuccess {
		response = TurnResponse{
			RoomName:  roomName,
			IsCreated: true,
		}
	} else {
		response = CreateRoomResponse{
			IsCreated:    false,
			RejectReason: "Room already created",
		}
	}

	data, _ := json.Marshal(response)

	message := Message{
		Type: "CreateRoomResponse",
		Raw:  data,
	}

	w.WriteJSON(message)
}
