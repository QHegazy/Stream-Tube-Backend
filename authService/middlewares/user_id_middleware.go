package middlewares

// import (
// 	"docs/internal/response"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func UserIdMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		cookie, err := c.Cookie("doc")
// 		if err != nil || cookie == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
// 				BaseResponse: response.BaseResponse{
// 					Status:  http.StatusUnauthorized,
// 					Message: "Unauthorized",
// 				},
// 				Error: "Session token is missing or invalid.",
// 			})
// 			return
// 		}
// 		c.Set("session_token", cookie)
// 		c.Next()
// 	}
// }
