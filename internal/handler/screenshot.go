package handler

import (
	"net/http"
	"thumburl-service/internal/service/imageservice"
	"thumburl-service/internal/service/screenshotservice"

	"github.com/gin-gonic/gin"
)

const GetScreenShotEndpoint = "/screenshot"

type GetScreenShotQuery struct {
	URL         string `form:"url" binding:"required"`
	Width       int    `form:"width" binding:"required"`
	Height      int    `form:"height" binding:"required"`
	ImageWidth  int    `form:"image_width" binding:"required"`
	ImageHeight int    `form:"image_height" binding:"required"`
}

func GetScreenShot(c *gin.Context) {
	var req GetScreenShotQuery
	if c.ShouldBindQuery(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid query"})
		return
	}

	data, err := screenshotservice.ScreenShot(req.URL, req.Width, req.Height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	data, err = imageservice.ResizeWebp(data, req.ImageWidth, req.ImageHeight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/webp", data)
}
