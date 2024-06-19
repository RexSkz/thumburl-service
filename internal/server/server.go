package server

import (
	"context"
	"os"
	"os/signal"

	"thumburl-service/internal/config"
	h "thumburl-service/internal/handler"
	"thumburl-service/internal/middleware"
	"thumburl-service/internal/pkg/logger"
	"thumburl-service/internal/service/screenshotservice"

	"github.com/gin-gonic/gin"
)

func Start() {
	config.Init()

	if err := screenshotservice.InitPool(); err != nil {
		logger.Panicw(
			context.Background(),
			"failed to initialize pool",
			"error", err,
		)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Infow(
			context.Background(),
			"shutting down",
		)
		screenshotservice.DisposePool()
		os.Exit(0)
	}()

	r := gin.Default()
	r.Use(
		middleware.Trace(),
		gin.RecoveryWithWriter(os.Stderr, func(c *gin.Context, err interface{}) {
			logger.Errorw(
				c,
				"recovery",
				"error", err,
			)
		}),
	)

	r.GET(h.GetScreenShotEndpoint, h.GetScreenShot)

	r.SetTrustedProxies(nil)
	r.Run(config.Config.Port)
}
