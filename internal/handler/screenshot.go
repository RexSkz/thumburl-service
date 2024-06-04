package handler

import (
	"net/http"
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
	var query GetScreenShotQuery
	if c.ShouldBindQuery(&query) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid query"})
		return
	}

	data, err := screenshotservice.ScreenShot(query.URL, query.Width, query.Height, query.ImageWidth, query.ImageHeight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/webp", data)
}
