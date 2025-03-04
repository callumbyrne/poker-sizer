package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/callumbyrne/poker-sizer/internal/services"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	roomService *services.RoomService
	upgrader    websocket.Upgrader
	clients     map[*websocket.Conn]string // map of connections to user IDs
}

func NewWebSocketHandler(roomService *services.RoomService) *WebSocketHandler {
	return &WebSocketHandler{
		roomService: roomService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections in dev
			},
		},
		clients: make(map[*websocket.Conn]string),
	}
}

type WSEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	// roomID := parts[3]
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	// Register the client
	h.clients[conn] = userID
	defer delete(h.clients, conn)

	// Main message handling loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var event WSEvent
		if err := json.Unmarshal(message, &event); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		// Handle the different event types
		switch event.Type {
		case "submit_vote":
			var payload struct {
				Value string `json:"value"`
			}
			json.Unmarshal(event.Payload, &payload)

			// TODO: Process the vote

		case "reveal_votes":
			// TODO: Reveal votes

		case "reset_voting":
			// TODO: Reset voting
		}
	}
}
