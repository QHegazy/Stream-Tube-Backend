package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors := os.Getenv("CORS")
		c.Header("Access-Control-Allow-Origin", cors)
		// c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Referer, User-Agent, DNT, Cache-Control, X-Mx-ReqToken, Keep-Alive, Connection, Host, Cookie, X-CSRF-Token, lg")
		// c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma, Date, Server, Authorization, Set-Cookie, lg")
		// c.Header("Access-Control-Allow-Credentials", "true")
		// c.Header("Access-Control-Max-Age", "86400")
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		// c.Header("X-Frame-Options", "DENY")
		// c.Header("X-Content-Type-Options", "nosniff")
		// c.Header("X-XSS-Protection", "1; mode=block")
		// c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		// c.Header("Referrer-Policy", "strict-origin")
		// c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		// c.Header("X-DNS-Prefetch-Control", "off")
		// c.Header("X-Download-Options", "noopen")
		// c.Header("X-Permitted-Cross-Domain-Policies", "none")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next() // Call the next handler
	}
}
