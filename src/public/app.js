new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        username: null, // Our username
        joined: false, // True if email and username have been filled in
        inRoom: false,
        rooms: [],
        newRoom: '',
        currentRoom: '',
        otherPlayerChoise: '',
        otherPlayerName: '',
        playerScore: 0,
        otherPlayerScore: 0,
        isRoomFull: false,
        isTurned: false,
        isOpponentTurned: false
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            if(msg.Type == "AuthResponse"){
                if(msg.Raw.isRegistred == true)
                {
                    self.username = $('<p>').html(this.username).text();
                    self.joined = true;
                    self.newRoom = $('<p>').html(this.newRoom).text()
                    var roomsArray = Object.keys(msg.Raw.roomsList).map(i => msg.Raw.roomsList[i])
                    Object.keys(msg.Raw.roomsList).forEach(key => {
                        var room = new Object;
                        room['name'] = key;
                        room['players'] = msg.Raw.roomsList[key];
                        self.rooms.push(room);
                    });
                    // var myData = Object.keys(msg.Raw.roomsList).map(key => {
                    //     return msg.Raw.roomsList[key];
                    // })
                    // self.rooms.concat(roomsArray)
                }
                else
                {
                    Materialize.toast(msg.Raw.rejectReason, 2000);
                }
            }
            else if(msg.Type == "CreateRoomResponse")
            {
                if(msg.Raw.isCreated == true)
                {
                    var keys = Object.keys( msg.Raw.roomName );
                    var key = keys[0]
                    var room = new Object;
                    room['name'] = key;
                    room['players'] = msg.Raw.roomName[key];
                    self.rooms.push(room);
                }
                else
                {
                    Materialize.toast(msg.Raw.rejectReason, 2000);
                }
            }
            else if(msg.Type == "EnterRoomResponse")
            {
                if(msg.Raw.isEntered == true)
                {
                    self.inRoom = true;
                    self.currentRoom = msg.Raw.roomname;
                }
                else
                {
                    Materialize.toast(msg.Raw.rejectReason, 2000);
                }
            }
            else if(msg.Type == "LeaveRoomResponse")
            {
                if(msg.Raw.isLeft == true)
                {
                    self.inRoom = false;
                    self.currentRoom = "";
                }
                else
                {
                    Materialize.toast(msg.Raw.rejectReason, 2000);
                }
            }
            else if(msg.Type == "TurnResponse")
            {
                if(msg.Raw.isApplied == true)
                {
                    if(msg.Raw.result != "")
                    {
                        Materialize.toast(msg.Raw.result, 2000);
                        self.otherPlayerScore = msg.Raw.otherPlayerScore;
                        self.playerScore = msg.Raw.currentPlayerScore;
                        self.isTurned = false;
                        self.isOpponentTurned = false;
                    } else{
                        if(msg.Raw.otherPlayerChoise == "OtherPlayerTurned"){
                            self.isOpponentTurned = true;
                        }else{
                            self.isTurned = true;
                        }
                    }
                }
                else
                {
                    Materialize.toast(msg.Raw.rejectReason, 2000);
                }
            }
            else if(msg.Type == "PlayerEneteredNotification")
            {
                self.isTurned = false;
                self.isRoomFull = true
                self.otherPlayerName = msg.Raw.otherPlayerName
            }
            else if(msg.Type == "PlayerLeftNotification")
            {
                self.isRoomFull = false
                self.otherPlayerName = ""
            }
        });
    },
    
    methods: {
        send: function () {
            if (this.newMsg != '') {
                var a = {
                    type: "message",
                    email: this.email,
                    username: this.username,
                    message: $('<p>').html(this.newMsg).text() // Strip out html
                }
                this.ws.send(
                    JSON.stringify(a));
                this.newMsg = ''; // Reset newMsg
            }
        },

        createRoom: function () {
            if (!this.newRoom) {
                Materialize.toast('You must enter room name', 2000);
                return
            }
            
            var createRequest = {
                roomName: this.newRoom
            }
            
            var wrappedCreateRequest = {
                type: "CreateRoomRequest",
                raw: createRequest
            }

            this.ws.send(JSON.stringify(wrappedCreateRequest));
        },

        join: function () {
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }

            var authRequest = {
                username: this.username
            }
            
            var wrappedAuthRequest = {
                type: "AuthRequest",
                raw: authRequest
            }

            this.ws.send(JSON.stringify(wrappedAuthRequest));
        },

        enterRoom:function (roomName){
            if (!roomName) {
                Materialize.toast('You must choose a room', 2000);
                return
            }

            var enterRoomRequest = {
                roomname: roomName
            }
            
            var wrappedEnterRoomRequest = {
                type: "EnterRoomRequest",
                raw: enterRoomRequest
            }

            this.ws.send(JSON.stringify(wrappedEnterRoomRequest));
        },

        leaveRoom:function (){
            if (!this.currentRoom) {
                Materialize.toast('You are not in room', 2000);
                return
            }

            var leaveRoomRequest = {
                roomname: this.currentRoom
            }
            
            var wrappedLeaveRoomRequest = {
                type: "LeaveRoomRequest",
                raw: leaveRoomRequest
            }

            this.ws.send(JSON.stringify(wrappedLeaveRoomRequest));
        },

        turn:function (playerChiose){
            if (!this.currentRoom) {
                Materialize.toast('You are not in room', 2000);
                return
            }

            var turnRequest = {
                choise: playerChiose
            }
            
            var wrappedturnRequest = {
                type: "TurnRequest",
                raw: turnRequest
            }

            this.ws.send(JSON.stringify(wrappedturnRequest));
        }
    }
});