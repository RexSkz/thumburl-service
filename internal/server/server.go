package server

import (
	h "thumburl-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	r.GET(h.GetScreenShotEndpoint, h.GetScreenShot)

	r.SetTrustedProxies(nil)
	r.Run(":8080")
}
