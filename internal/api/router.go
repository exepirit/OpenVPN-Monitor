package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request, server *Server) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed ", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Status request from %s", r.RemoteAddr)
	status, err := server.RequestStatus()
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		log.Print("[ERR] ", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
