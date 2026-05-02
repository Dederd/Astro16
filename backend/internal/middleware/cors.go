package middleware

import (
	"strings"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	allowedOrigins := []string{
		"https://astro16-jolt.vercel.app",
		"http://localhost:5173",
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow semua subdomain vercel.app atau origin yang ada di list
		allowed := false
		for _, o := range allowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}
		// Allow semua preview URL Vercel
		if strings.Contains(origin, "vercel.app") {
			allowed = true
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Session-ID, X-Admin-Key")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
