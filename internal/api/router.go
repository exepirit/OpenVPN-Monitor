package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/exepirit/OpenVPN-Monitor/internal/openvpn"
)

func StatusHandler(serverAddr string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		server := openvpn.Server{
			Address: serverAddr,
		}
		if err := server.Connect(); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer server.Close()

		c, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		status, err := server.RequestStatus(c)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, status)
	}
}
