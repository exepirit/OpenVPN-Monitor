package main

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
	"os"

	"github.com/exepirit/OpenVPN-Monitor/internal/api"
	"github.com/exepirit/OpenVPN-Monitor/internal/openvpn"
)

func HandleHTTP(address string, server *openvpn.Server) error {
	httpSrv := gin.Default()
	httpSrv.GET("/api/status", api.StatusHandler(server))
	httpSrv.StaticFile("/", "static/index.html")
	return httpSrv.Run(address)
}

func appFunc(ctx *cli.Context) error {
	server := openvpn.Server{Address: ctx.String("server")}
	if err := server.Connect(); err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	return HandleHTTP(ctx.String("bind"), &server)
}

func main() {
	app := &cli.App{
		Name: "openvpn-monitor",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server",
				Aliases:     []string{"s"},
				EnvVars:     []string{"OPENVPN_SERVER"},
				Required:    true,
				DefaultText: "localhost:7505",
			},
			&cli.StringFlag{
				Name:        "bind",
				Aliases:     []string{"b"},
				DefaultText: "0.0.0.0:8000",
			},
		},
		Action: appFunc,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
