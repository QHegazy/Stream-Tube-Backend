package middlewares

// import (
// 	"docs/internal/response"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func BadRequestMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Check if a required query parameter is present
// 		if c.Query("requiredParam") == "" {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
// 				BaseResponse: response.BaseResponse{
// 					Status:  http.StatusBadRequest,
// 					Message: "Bad Request",
// 				},
// 				Error: "Required query parameter 'requiredParam' is missing.",
// 			})
// 			return
// 		}

// 		// Check if the JSON body is empty for POST requests
// 		if c.Request.Method == http.MethodPost {
// 			var json gin.H
// 			if err := c.ShouldBindJSON(&json); err != nil || len(json) == 0 {
// 				c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
// 					BaseResponse: response.BaseResponse{
// 						Status:  http.StatusBadRequest,
// 						Message: "Bad Request",
// 					},
// 					Error: "JSON body is empty",
// 				})
// 				return
// 			}
// 		}

// 		c.Next()
// 	}
// }
