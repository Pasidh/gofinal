package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	fmt.Println("start #middleware")

	token := c.GetHeader("Authorization")
	if token != "Bearer token1234" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you don't have the permission!!"})
		c.Abort()
		return
	}

	c.Next()

	fmt.Println("end #middleware")
}
