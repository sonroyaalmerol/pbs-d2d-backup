//go:build windows

package controllers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/sonroyaalmerol/pbs-plus/internal/agent/snapshots"
)

func JobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusBadRequest)
	}

	driveLetter := r.PathValue("drive")
	if driveLetter == "" {
		http.Error(w, "Drive not found", http.StatusNotFound)
	}

	snapshot, err := snapshots.Snapshot(driveLetter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	port := 33451
	for {
		if port >= 33476 {
			port = 33451
		}
		listenAt := fmt.Sprintf("0.0.0.0:%d", port)
		listener, err := net.Listen("tcp", listenAt)
		if err == nil {
			listener.Close()
			break
		}
		listener.Close()
		port++
	}

	go func() {
		err := snapshot.Serve(port)
		if err != nil {

		}
	}()

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"port": port})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
