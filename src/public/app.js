new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        username: null, // Our username
        joined: false, // True if email and username have been filled in
        rooms: '',
        newRoom: ''
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
                    self.rooms += '<li class="collection-item"><div>'+msg.Raw.roomName+'<a href="room" class="secondary-content"><i class="material-icons">meeting_room</i></a></div></li>'
                }
                else
                {
                    Materialize.toast(msg.Raw.rejectReason, 2000);
                }

            }
            // if (msg.type == 'message') {
            //     self.chatContent += '<div class="chip">'
            //     + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
            //     + msg.username
            //     + '</div>'
            //     + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
            //     var element = document.getElementById('chat-messages');
            //     element.scrollTop = element.scrollHeight;
            // } else if(msg.type == 'createRoom'){
            //     //self.rooms += '<div class="chip">'+ 'Hello'+'</div>'
            //     self.rooms += '<li class="collection-item"><div>'+msg.message+'<a href="room" class="secondary-content"><i class="material-icons">meeting_room</i></a></div></li>'
    
            //     var element = document.getElementById('rooms-list');
            //     element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
            // }
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
    }
});