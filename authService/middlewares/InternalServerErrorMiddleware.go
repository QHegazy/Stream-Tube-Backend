package middlewares

// import (
// 	"docs/internal/response"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// // InternalServerErrorMiddleware recovers from panics and returns an internal server error response
// func InternalServerErrorMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				log.Printf("Internal Server Error: %v", err)
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{
// 					BaseResponse: response.BaseResponse{
// 						Status:  http.StatusInternalServerError,
// 						Message: "Internal Server Error",
// 					},
// 					Error: "An unexpected error occurred.",
// 				})
// 			}
// 		}()
// 		c.Next()
// 	}
// }
