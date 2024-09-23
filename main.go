package main

import (
	"gowebsocket-lib/websocket"
	"log"
	"net/http"
	"time"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Perform WebSocket upgrade
	err := websocket.Upgrade(w, r, []string{"chat", "binary", "json"})
	if err != nil {
		http.Error(w, "Could not upgrade", http.StatusBadRequest)
		return
	}

	// Hijack the connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	conn, _, err := hijacker.Hijack() // Hijack() returns net.Conn, *bufio.ReadWriter, error
	if err != nil {
		log.Println("Error hijacking connection:", err)
		return
	}
	defer conn.Close()

	// Create a new WebSocket connection
	wsConn := websocket.NewConn(conn)

	// Start a goroutine to send periodic messages
	ticker := time.NewTicker(5 * time.Second) // Change the duration as needed
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				message := []byte("Periodic message from server")
				if err := wsConn.WriteMessage(message); err != nil {
					log.Println("Write error:", err)
					return
				}
			}
		}
	}()

	// Handle WebSocket communication
	for {
		messageType, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		if messageType == "binary" {
			log.Println("Received message:", msg)

		} else {
			log.Println("Received message:", string(msg))
		}

		// Echo message back
		if err := wsConn.WriteMessage(msg); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
