package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com/exepirit/OpenVPN-Monitor/internal/api"
	"github.com/exepirit/OpenVPN-Monitor/internal/openvpn"
)

type Configuration struct {
	OpenVPNAddress string `json:"openvpn_address"`
	HandleAddress  string `json:"http_address"`
}

func HandleHTTP(address string, server *openvpn.Server) {
	httpSrv := gin.Default()
	httpSrv.GET("/api/status", api.StatusHandler(server))
	httpSrv.StaticFile("/", "static/index.html")
	if err := httpSrv.Run(address); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := Configuration{}
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(cfgFile).Decode(&config); err != nil {
		log.Fatal(err)
	}

	server := openvpn.Server{Address: config.OpenVPNAddress}
	if err := server.Connect(); err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	HandleHTTP(config.HandleAddress, &server)
}
