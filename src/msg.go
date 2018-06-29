package main

import (
	"encoding/json"
)

// Define our message object
type Message struct {
	Type string
	Raw  json.RawMessage
}

// Define our message object
type AuthRequest struct {
	Username string `json:"username"`
}

// Define our message object
type AuthResponse struct {
	IsRegistred  bool   `json:"isRegistred"`
	RejectReason string `json:"rejectReason"`
}

// Define our message object
type CreateRoomRequest struct {
	RoomName string `json:"roomname"`
}

// Define our message object
type CreateRoomResponse struct {
	IsCreated    bool   `json:"isCreated"`
	RoomName     string `json:"roomName"`
	RejectReason string `json:"rejectReason"`
}

type TurnRequest struct {
	Choise string `json:"choise"`
}

type TurnResponse struct {
	IsApplied    bool   `json:"isApplied"`
	Result       string `json:"result"`
	RejectReason string `json:"rejectReason"`
}

// Define our message object
type EnterRoomRequest struct {
	RoomName string `json:"roomname"`
}

// Define our message object
type EnterRoomResponse struct {
	IsEntered    bool   `json:"isEntered"`
	RoomName     string `json:"roomname"`
	RejectReason string `json:"rejectReason"`
}
