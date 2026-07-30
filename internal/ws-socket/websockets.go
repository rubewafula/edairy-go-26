package sockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var allowedOrigins = map[string]bool{
	"https://arithi.edairy.africa":           true,
	"https://api.arithi.edairy.africa":       true,
	"https://tigania-west.edairy.africa":     true,
	"https://api.tigania-west.edairy.africa": true,
	"https://mukululu.edairy.africa":         true,
	"https://api.mukululu.edairy.africa":     true,
	"https://mwimbi.edairy.africa":           true,
	"https://api.mwimbi.edairy.africa":       true,
	"https://mutuati.edairy.africa":          true,
	"https://api.mutuati.edairy.africa":      true,
	"https://nkuene.edairy.africa":           true,
	"https://api.nkuene.edairy.africa":       true,
	"https://dev.edairy.africa":              true,
	"https://api.dev.edairy.africa":          true,
	"https://edairy.africa":                  true,
	"http://localhost:5173":                  true,
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		return allowedOrigins[origin]
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

func validateTokenAndGetUserID(tokenStr string, jwtSecret []byte) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid user id")
	}

	return strconv.FormatUint(uint64(userID), 10), nil
}

func extractToken(c *gin.Context) string {
	if token := c.Query("token"); token != "" {
		return token
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
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
					log.Printf("Buffer full for user %s, skipping broadcast", client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

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

func ServeWS(hub *Hub, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			log.Println("WebSocket upgrade failed: token missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
			return
		}

		userID, err := validateTokenAndGetUserID(tokenStr, jwtSecret)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
