package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/thevedantmod/stand-clear/board"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	/* change here if need to restrict API access. this is public data now, so will allow anyone */
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")

	line := r.URL.Query().Get("line")
	stopID := r.URL.Query().Get("stop_id")
	nStr := r.URL.Query().Get("N")

	if nStr == "" {
		nStr = "10"
	}

	N, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, "Invalid N parameter", http.StatusBadRequest)
		return
	}

	fmt.Printf("line is '%s'", line)

	arrivals, err := board.GetArrivals(line, stopID, N)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(arrivals)
}
