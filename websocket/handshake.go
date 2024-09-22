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
func Upgrade(w http.ResponseWriter, r *http.Request) error {
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

	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", acceptKey)
	w.WriteHeader(http.StatusSwitchingProtocols)

	return nil
}
