package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// AcceptKey computes the Sec-WebSocket-Accept value
func AcceptKey(key string) string {
	const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.New()
	h.Write([]byte(key + wsGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Upgrade upgrades the HTTP connection to a WebSocket connection
func Upgrade(w http.ResponseWriter, r *http.Request, supportedProtocols []string) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("method not allowed")
	}

	if !strings.Contains(r.Header.Get("Connection"), "Upgrade") ||
		r.Header.Get("Upgrade") != "websocket" {
		return fmt.Errorf("bad upgrade")
	}

	wsKey := r.Header.Get("Sec-WebSocket-Key")
	if wsKey == "" {
		return fmt.Errorf("missing Sec-WebSocket-Key")
	}

	acceptKey := AcceptKey(wsKey)

	// Handle Sec-WebSocket-Protocol negotiation
	clientProtocols := r.Header.Get("Sec-WebSocket-Protocol")
	var selectedProtocol string
	if clientProtocols != "" {
		clientProtocolList := strings.Split(clientProtocols, ",")
		// Check if any client requested protocol is supported by the server
		for _, clientProtocol := range clientProtocolList {
			clientProtocol = strings.TrimSpace(clientProtocol)
			for _, serverProtocol := range supportedProtocols {
				if clientProtocol == serverProtocol {
					selectedProtocol = clientProtocol
					break
				}
			}
			if selectedProtocol != "" {
				break
			}
		}
	}

	// Set necessary headers for WebSocket upgrade
	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", acceptKey)

	// If a protocol is selected, set the Sec-WebSocket-Protocol header
	if selectedProtocol != "" {
		w.Header().Set("Sec-WebSocket-Protocol", selectedProtocol)
		fmt.Println("Selected Protocol:", selectedProtocol)
	} else {
		fmt.Println("No protocol selected")
	}

	w.WriteHeader(http.StatusSwitchingProtocols)
	return nil
}
