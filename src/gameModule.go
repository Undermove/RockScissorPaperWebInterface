package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type GameModule struct {
	loosersMap   map[string]string
	authModule   *AuthModule
	roomsManager *RoomsManager
}

func NewGameModule(authModule *AuthModule, roomsManager *RoomsManager) *GameModule {

	var loosersMap = map[string]string{
		"Rock":     "Scissors",
		"Paper":    "Rock",
		"Scissors": "Paper",
	}

	return &GameModule{
		loosersMap:   loosersMap,
		authModule:   authModule,
		roomsManager: roomsManager,
	}
}

func (g *GameModule) turn(playerOne *Player, playerTwo *Player) string {
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

func (g *GameModule) Turn(w *websocket.Conn, request TurnRequest) *TurnResponse {
	room := g.roomsManager.ConnToRooms[w]

	currentPlayerName := g.authModule.Clients[w]
	ok, currentPlayer := room.TryGetCurrentPlayer(currentPlayerName)
	if !ok {
		return &TurnResponse{
			RejectReason: "You must be in room first",
			IsApplied:    false,
		}
	}
	currentPlayer.Choise = request.Choise

	if ok, otherPlayer := room.TryGetOtherPlayer(g.authModule.Clients[w]); ok {
		if otherPlayer.Choise != "" {
			winner := gameModule.turn(currentPlayer, otherPlayer)
			return &TurnResponse{
				Result:             winner,
				OtherPlayerChoise:  otherPlayer.Choise,
				IsApplied:          true,
				CurrentPlayerScore: currentPlayer.Score,
				OtherPlayerScore:   otherPlayer.Score,
			}
		}
	}

	return &TurnResponse{
		IsApplied: true,
	}
}

func (g *GameModule) SendTurnResponse(resp TurnResponse, w *websocket.Conn) {
	data, _ := json.Marshal(resp)
	var forOtherPlayer = TurnResponse{
		Result:             resp.Result,
		OtherPlayerChoise:  "OtherPlayerTurned",
		IsApplied:          true,
		CurrentPlayerScore: resp.OtherPlayerScore,
		OtherPlayerScore:   resp.CurrentPlayerScore,
	}
	dataForOther, _ := json.Marshal(forOtherPlayer)

	message := Message{
		Type: "TurnResponse",
		Raw:  data,
	}

	messageForOther := Message{
		Type: "TurnResponse",
		Raw:  dataForOther,
	}

	room := g.roomsManager.ConnToRooms[w]
	_, otherPlayer := room.TryGetOtherPlayer(g.authModule.Clients[w])

	// send message to all players
	g.authModule.AuthClients[otherPlayer.Name].WriteJSON(messageForOther)
	w.WriteJSON(message)
	return
}
