package middlewares

// import (
// 	"docs/internal/models"
// 	"docs/internal/response"
// 	"docs/internal/services/sessions"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		cookie, err := c.Cookie("lg")
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

// func CheckSessionToken() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		res := make(chan models.Session, 1)
// 		sessionToken, exists := c.Get("session_token")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
// 				BaseResponse: response.BaseResponse{
// 					Status:  http.StatusUnauthorized,
// 					Message: "Unauthorized",
// 				},
// 				Error: "Session token not found.",
// 			})
// 			return
// 		}

// 		token := sessionToken.(string)
// 		go sessions.GetSession(token, res)

// 		select {
// 		case session := <-res:
// 			currentTime := time.Now().UTC()
// 			if session == (models.Session{}) || currentTime.After(session.ExpiresAt) {

// 				c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
// 					BaseResponse: response.BaseResponse{
// 						Status:  http.StatusUnauthorized,
// 						Message: "Unauthorized",
// 					},
// 					Error: "Session token is missing or invalid.",
// 				})
// 				return
// 			}
// 		case <-time.After(time.Second * 5):
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{
// 				BaseResponse: response.BaseResponse{
// 					Status:  http.StatusInternalServerError,
// 					Message: "Unauthorized",
// 				},
// 				Error: "Timeout while retrieving session.",
// 			})
// 			return
// 		}

// 		c.Next()
// 	}
// }
