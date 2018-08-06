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
        otherPlayerChoise: ''
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
                    self.rooms = msg.Raw.roomsList
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
                    self.rooms.push(msg.Raw.roomName);
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
                    self.currentRoom = msg.Raw.roomname
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
                    self.currentRoom = ""
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
                    }
                }
                else
                {
                    Materialize.toast(msg.Raw.rejectReason, 2000);
                }
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