package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":    "success!",
		"message":   message,
		"data":      data,
		"timestamp": time.Now().Unix(),
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":    "error!",
		"message":   message,
		"timestamp": time.Now().Unix(),
	})
}
