package main

import (
	"OpenVPN-Monitor/api"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Configuration struct {
	OpenVPNAddress string `json:"openvpn_address"`
	HandleAddress  string `json:"http_address"`
}

func HandleHTTP(address string, server *api.Server) {
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         address,
	}

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		api.Status(w, r, server)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	log.Print("Handling ", address)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func main() {
	config := Configuration{}
	cfgFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(cfgFile).Decode(&config); err != nil {
		panic(err)
	}
	server := api.Server{Address: config.OpenVPNAddress}
	if err := server.Connect(); err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	HandleHTTP(config.HandleAddress, &server)
}
