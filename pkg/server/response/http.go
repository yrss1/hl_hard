package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Data":    data,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{
		"Success": true,
		"Data":    data,
	})
}

func BadRequest(c *gin.Context, err error, data any) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Success": false,
		"Message": err.Error(),
		"Data":    data,
	})
}

func NotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, gin.H{
		"Success": false,
		"Message": err.Error(),
	})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"Success": false,
		"Message": err.Error(),
	})
}

func MethodNotAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedMethods := map[string]bool{
			http.MethodGet:    true,
			http.MethodPost:   true,
			http.MethodPut:    true,
			http.MethodDelete: true,
		}

		if !allowedMethods[c.Request.Method] {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method Not Allowed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
