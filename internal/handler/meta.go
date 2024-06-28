package handler

import (
	"net/http"

	"thumburl-service/internal/service/metaservice"

	"github.com/gin-gonic/gin"
)

const GetMetaEndpoint = "/meta"

type GetMetaQuery struct {
	URL string `form:"url" binding:"required"`
}

func GetMeta(c *gin.Context) {
	var req GetMetaQuery
	if c.ShouldBindQuery(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid query"})
		return
	}

	data, err := metaservice.GetMeta(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
