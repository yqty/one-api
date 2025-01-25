package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// API 路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "One API is running",
			"version": "1.0.0",
		})
	})

	router.ServeHTTP(w, r)
}
