package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/callumbyrne/poker-sizer/internal/services"
	"github.com/callumbyrne/poker-sizer/web/templates/pages"
)

type RoomHandler struct {
	roomService *services.RoomService
	templates   *template.Template
}

func NewRoomHandler(roomService *services.RoomService) *RoomHandler {
	templates := template.Must(template.ParseGlob("web/templates/*.html"))
	template.Must(templates.ParseGlob("web/templates/components/*.html"))

	return &RoomHandler{
		roomService: roomService,
		templates:   templates,
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	pages.Home().Render(r.Context(), w)
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	roomName := r.Form.Get("name")
	room, err := h.roomService.CreateRoom(roomName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		http.Redirect(w, r, "/rooms/"+room.ID, http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": room.ID})
	}
}

func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}

	roomID := parts[2]
	room, err := h.roomService.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	h.templates.ExecuteTemplate(w, "room.html", room)
}
