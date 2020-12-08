package api

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/exepirit/OpenVPN-Monitor/internal/openvpn"
)

func StatusHandler(server *openvpn.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status, err := server.RequestStatus()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, status)
	}
}
