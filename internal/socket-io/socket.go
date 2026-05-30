package sockets

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

// Server is a global reference to the Socket.IO server
var Server *socketio.Server

// Init initializes the Socket.IO server and sets up event handlers
func Init() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("Socket connected:", s.ID())
		return nil
	})

	// Client should emit 'join' with their user ID after connecting
	server.OnEvent("/", "join", func(s socketio.Conn, userID string) {
		if userID != "" {
			s.Join("user_" + userID)
			log.Printf("Socket %s joined room: user_%s", s.ID(), userID)
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("Socket error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("Socket %s disconnected: %s", s.ID(), reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO serve error: %v", err)
		}
	}()

	Server = server
}

// EmitNotification sends a real-time notification to a specific user's room
func EmitNotification(userID string, notification interface{}) {
	if Server != nil {
		Server.BroadcastToRoom("/", "user."+userID, "notification", notification)
	}
}
