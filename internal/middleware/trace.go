package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIDKey = "trace_id"
const TraceIDGlobal = "global"

func Trace() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set(TraceIDKey, traceID)
		c.Next()
	}
}
