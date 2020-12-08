package main

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
	"os"

	"github.com/exepirit/OpenVPN-Monitor/internal/api"
)

func HandleHTTP(address string, serverAddr string) error {
	httpSrv := gin.Default()
	httpSrv.GET("/api/status", api.StatusHandler(serverAddr))
	httpSrv.StaticFile("/", "static/index.html")
	return httpSrv.Run(address)
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
		Action: func(ctx *cli.Context) error {
			return HandleHTTP(ctx.String("bind"), ctx.String("server"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
