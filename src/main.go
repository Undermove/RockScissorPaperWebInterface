package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

var roomsManager *RoomsManager
var gameModule *GameModule
var authModule *AuthModule                  // handle clients connections
var broadcast = make(chan WebSocketMessage) // broadcast channel
var logpath = flag.String("logpath", os.Getenv("RPS_LOG_PATH"), "Log Path")

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketMessage struct {
	Message Message
	fromWs  *websocket.Conn
}

func provideScriptFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, os.Getenv("RPS_PATH")+"/public/app.js")
}

func provideStyleFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, os.Getenv("RPS_PATH")+"/public/style.css")
}

var dir string

func main() {
	NewLog(*logpath)
	dir = os.Getenv("RPS_PATH")
	Log.Println(os.Getenv("RPS_PATH"))
	authModule = NewModule()
	roomsManager = NewRoomsManager(authModule)
	gameModule = NewGameModule(authModule, roomsManager)

	// Create a simple file server
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(os.Getenv("RPS_PATH") + "/public"))
	r.Handle("/", http.StripPrefix("/", fs))
	r.HandleFunc("/app.js", provideScriptFile).Methods("GET")
	r.HandleFunc("/style.css", provideStyleFile).Methods("GET")

	// Configure websocket route
	r.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	Log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		Log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	authModule.AddConnection(ws)

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			Log.Printf("error: %v", err)
			authModule.Disconnect(ws)
			break
		}
		// Send the newly received message to the broadcast channel
		wsMsg := WebSocketMessage{fromWs: ws, Message: msg}
		broadcast <- wsMsg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		wsMsg := <-broadcast

		if authModule.IsLoggedIn(wsMsg.fromWs) {
			if wsMsg.Message.Type == "CreateRoomRequest" {
				processCreateRoomRequest(wsMsg)
			} else if wsMsg.Message.Type == "TurnRequest" {
				processTurnRequest(wsMsg)
			} else if wsMsg.Message.Type == "EnterRoomRequest" {
				processEnterRoomRequest(wsMsg)
			} else if wsMsg.Message.Type == "LeaveRoomRequest" {
				processLeaveRoomRequest(wsMsg)
			}
		} else if wsMsg.Message.Type == "AuthRequest" {
			processAuthRequest(wsMsg)
		}
	}
}

func processAuthRequest(wsMsg WebSocketMessage) {
	var request AuthRequest
	err := json.Unmarshal(wsMsg.Message.Raw, &request)
	if err != nil {
		return
	}
	if authModule.Authenticate(wsMsg.fromWs, request) {
		authModule.SendSuccessResponse(wsMsg.fromWs, roomsManager.GetRoomStats())
	} else {
		authModule.SendRejectResponse(wsMsg.fromWs)
	}
}

func processCreateRoomRequest(wsMsg WebSocketMessage) {
	var request CreateRoomRequest

	err := json.Unmarshal(wsMsg.Message.Raw, &request)
	if err != nil {
		return
	}

	if roomsManager.AddRoom(request.RoomName) {
		roomsManager.SendResponse(true, request.RoomName, wsMsg.fromWs)
	} else {
		roomsManager.SendResponse(false, request.RoomName, wsMsg.fromWs)
	}
}

func processEnterRoomRequest(wsMsg WebSocketMessage) {
	var request EnterRoomRequest

	err := json.Unmarshal(wsMsg.Message.Raw, &request)
	if err != nil {
		return
	}

	if roomsManager.EnterRoom(wsMsg.fromWs, request.RoomName) {
		roomsManager.SendRoomEnterResponse(true, request.RoomName, wsMsg.fromWs)
		roomsManager.SendPlayerEnteredNotification(request.RoomName, wsMsg.fromWs)
	} else {
		roomsManager.SendRoomEnterResponse(false, request.RoomName, wsMsg.fromWs)
	}
}

func processLeaveRoomRequest(wsMsg WebSocketMessage) {
	var request LeaveRoomRequest

	err := json.Unmarshal(wsMsg.Message.Raw, &request)
	if err != nil {
		return
	}

	if roomsManager.LeaveRoom(wsMsg.fromWs, authModule.Clients[wsMsg.fromWs]) {
		roomsManager.SendLeaveRoomResponse(true, request.RoomName, wsMsg.fromWs)
		roomsManager.SendPlayerLeftNotification(request.RoomName, wsMsg.fromWs)
	} else {
		roomsManager.SendLeaveRoomResponse(false, request.RoomName, wsMsg.fromWs)
	}
}

func processTurnRequest(wsMsg WebSocketMessage) {
	var request TurnRequest

	err := json.Unmarshal(wsMsg.Message.Raw, &request)
	if err != nil {
		return
	}
	turnResp := gameModule.Turn(wsMsg.fromWs, request)
	gameModule.SendTurnResponse(*turnResp, wsMsg.fromWs)
}
