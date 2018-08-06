package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Player struct {
	Name   string
	Score  int
	Choise string
}

// Room with players
type Room struct {
	Name    string
	Players [2]*Player
}

func (r *Room) TryGetOtherPlayer(currentPlayerName string) (bool, *Player) {
	for i := 0; i < 2; i++ {
		if r.Players[i] != nil && r.Players[i].Name != currentPlayerName {
			return true, r.Players[i]
		}
	}

	return false, nil
}

func (r *Room) TryGetCurrentPlayer(currentPlayerName string) (bool, *Player) {
	for i := 0; i < 2; i++ {
		if r.Players[i] != nil && r.Players[i].Name == currentPlayerName {
			return true, r.Players[i]
		}
	}

	return false, nil
}

func (r *Room) EnterRoom(player Player) bool {
	if len(r.Players) <= 2 {
		for i := 0; i < 2; i++ {
			if r.Players[i] == nil {
				r.Players[i] = &player
				return true
			}
		}
	}

	return false
}

func (r *Room) LeaveRoom(playerName string) bool {
	for i := 0; i < 2; i++ {
		if r.Players[i] != nil && r.Players[i].Name == playerName {
			r.Players[i] = nil
			return true
		}
	}

	return false
}

type RoomsManager struct {
	Rooms       map[string]*Room
	ConnToRooms map[*websocket.Conn]*Room
	AuthModule  *AuthModule
}

func NewRoomsManager(auth *AuthModule) *RoomsManager {
	return &RoomsManager{
		Rooms:       make(map[string]*Room),
		ConnToRooms: make(map[*websocket.Conn]*Room),
		AuthModule:  auth,
	}
}

func (rm *RoomsManager) GetRoomNames() []string {
	// Can this be done better?
	roomsList := make([]string, len(rm.Rooms))
	idx := 0
	for _, value := range rm.Rooms {
		roomsList[idx] = value.Name
		idx++
	}

	return roomsList
}

func (rm *RoomsManager) SendResponse(isSuccess bool, roomName string, w *websocket.Conn) {
	var response CreateRoomResponse

	if isSuccess {
		response = CreateRoomResponse{
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

func (rm *RoomsManager) SendRoomEnterResponse(isSuccess bool, roomName string, w *websocket.Conn) {
	var response EnterRoomResponse

	if isSuccess {
		response = EnterRoomResponse{
			RoomName:  roomName,
			IsEntered: true,
		}
	} else {
		response = EnterRoomResponse{
			IsEntered:    false,
			RejectReason: "Room is full",
		}
	}

	data, _ := json.Marshal(response)

	message := Message{
		Type: "EnterRoomResponse",
		Raw:  data,
	}

	w.WriteJSON(message)
}

func (rm *RoomsManager) AddRoom(roomName string) bool {
	if _, ok := rm.Rooms[roomName]; ok {
		return false
	}
	rm.Rooms[roomName] = &Room{
		Name:    roomName,
		Players: [2]*Player{nil, nil},
	}
	return true
}

func (rm *RoomsManager) EnterRoom(ws *websocket.Conn, roomName string) bool {
	isEneterd := rm.Rooms[roomName].EnterRoom(Player{
		Name:   rm.AuthModule.Clients[ws],
		Choise: "",
		Score:  0,
	})

	if isEneterd {
		rm.ConnToRooms[ws] = rm.Rooms[roomName]
	}

	return isEneterd
}

func (rm *RoomsManager) LeaveRoom(ws *websocket.Conn, username string) bool {
	if _, ok := rm.ConnToRooms[ws]; ok {
		rm.ConnToRooms[ws].LeaveRoom(username)
		delete(rm.ConnToRooms, ws)
	}

	return true
}

func (rm *RoomsManager) IsInRoom(ws *websocket.Conn) bool {
	if _, ok := rm.ConnToRooms[ws]; ok {
		return true
	}

	return false
}

func (rm *RoomsManager) SendLeaveRoomResponse(isSuccess bool, roomName string, w *websocket.Conn) {
	var response LeaveRoomResponse

	if isSuccess {
		response = LeaveRoomResponse{
			RoomName: roomName,
			IsLeft:   true,
		}
	} else {
		response = LeaveRoomResponse{
			IsLeft:       false,
			RejectReason: "You are not in room",
		}
	}

	data, _ := json.Marshal(response)

	message := Message{
		Type: "LeaveRoomResponse",
		Raw:  data,
	}

	w.WriteJSON(message)
}
