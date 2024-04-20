package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Event represents the structure of our JSON payload
type Event struct {
	Data string `json:"data"`
}

func main() {
	http.HandleFunc("/log-event", logEvent)
	fmt.Println("Server is listening on port 8090...")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}

func logEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the event to JSON for logging
	eventJSON, err := json.Marshal(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON payload to the events.log file
	file, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if _, err := file.Write(eventJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := file.WriteString("\n"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response back indicating success
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("id-token", "GUID:ADSF-23de-324D-ffff")

	response := map[string]string{"message": "Event logged successfully"}
	json.NewEncoder(w).Encode(response)
}
