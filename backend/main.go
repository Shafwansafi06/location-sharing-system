package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// upgrader configures the WebSocket upgrading process.
// It specifies buffer sizes and an origin check.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin determines if a connection request is allowed.
	// We return true globally here to easily allow local React dev servers to connect.
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

// serveWs handles incoming HTTP websocket setup requests from the frontend.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 1. Parse the groupID from the URL path: /ws/{groupID}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	groupID := pathParts[2]

	// 2. Extract the userID from the query string
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "Missing userID query parameter", http.StatusBadRequest)
		return
	}

	// 3. Upgrade the standard HTTP connection to a persistent WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	// 4. Create a new Client instance that wrapping the WebSocket.
	client := &Client{
		hub:     hub,
		conn:    conn,
		// Buffered channel to queue outgoing messages safely (avoids blocking).
		send:    make(chan []byte, 256), 
		GroupID: groupID,
		UserID:  userID,
	}

	// 5. Tell the Hub that a new client has joined.
	hub.register <- client

	// 6. Launch the background Goroutines to pump data in and out independently.
	go client.WritePump()
	go client.ReadPump()
}

// corsMiddleware intercepts HTTP requests to explicitly attach Cross-Origin headers.
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Instantly approve OPTIONS pre-flight checks strictly required by Web Browsers.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Move on to the actual handler logic
		next(w, r)
	}
}

func main() {
	// Initialize the centralized Hub which synchronizes group maps seamlessly.
	hub := NewHub()
	
	// Start the Hub operating on its own separate Goroutine loop.
	go hub.Run()

	// Direct endpoint configuring real-time peer groups.
	http.HandleFunc("/ws/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	}))
	
	// A basic health check to instantly verify if our underlying Go server is alive.
	http.HandleFunc("/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	log.Println("Backend listening on :8080")
	// Start local HTTP server listening continuously on port 8080.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
