package sockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for the websocket handshake
	},
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

var Manager *Hub
var once sync.Once

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func InitHub(hub *Hub) {
	once.Do(func() {
		Manager = hub
	})
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("User registered: %s", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					// If the buffer is full, the unregister process will handle cleanup
					log.Printf("Buffer full for user %s, skipping broadcast", client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Join associates a client with a specific user ID and registers them with the hub
func (h *Hub) Join(userID string, client *Client) {
	client.ID = userID
	h.Register <- client
}

func readPump(hub *Hub, c *Client) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// handle incoming messages
		hub.Broadcast <- message
	}
}

func writePump(c *Client) {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

// ServeWS upgrades the HTTP connection to a WebSocket and registers the client
func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 👉 THIS is your gin context
		userID := c.GetString("user_id") // from auth middleware
		if userID == "" {
			userID = c.Query("user_id") // fallback (not secure)
		}
		if userID == "" {
			log.Println("WebSocket upgrade failed: userID query parameter missing")
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &Client{
			ID:   userID,
			Conn: conn,
			Send: make(chan []byte, 256),
		}

		hub.Join(userID, client)

		go readPump(hub, client)
		go writePump(client)
	}
}

// EmitNotification sends a real-time notification to a specific user's room
// EmitNotification sends a real-time notification to a specific user
func EmitNotification(userID string, notification interface{}) {

	if Manager == nil {
		return
	}

	Manager.mu.RLock()
	client, ok := Manager.Clients[userID]
	Manager.mu.RUnlock()

	if !ok {
		return
	}

	data, err := json.Marshal(notification)
	if err != nil {
		log.Printf("failed to marshal notification: %v", err)
		return
	}

	select {
	case client.Send <- data:
	default:
		log.Printf("user_%s channel full, dropping message", userID)
	}
}
