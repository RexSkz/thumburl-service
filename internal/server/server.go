package server

import (
	"fmt"
	"os"
	"os/signal"
	h "thumburl-service/internal/handler"
	"thumburl-service/internal/service/screenshotservice"

	"github.com/gin-gonic/gin"
)

func Start() {
	if err := screenshotservice.InitPool(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("Shutting down...\n")
		screenshotservice.Dispose()
		os.Exit(0)
	}()

	r := gin.Default()

	r.GET(h.GetScreenShotEndpoint, h.GetScreenShot)

	r.SetTrustedProxies(nil)
	r.Run(":8080")
}
