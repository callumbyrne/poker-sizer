package main

import (
	"log"
	"net/http"
	"os"

	"github.com/callumbyrne/poker-sizer/internal/handlers"
	"github.com/callumbyrne/poker-sizer/internal/services"
	"github.com/callumbyrne/poker-sizer/internal/store"
)

func main() {
	roomStore := store.NewMemoryStore()

	roomService := services.NewRoomService(roomStore)

	roomHandler := handlers.NewRoomHandler(roomService)
	wsHandler := handlers.NewWebSocketHandler(roomService)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/rooms/create", roomHandler.CreateRoom)
	// mux.HandleFunc("/rooms/join", roomHandler.JoinRoom)
	mux.HandleFunc("/rooms/{id}", roomHandler.GetRoom)
	mux.HandleFunc("/ws/rooms/{id}", wsHandler.HandleConnection)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
