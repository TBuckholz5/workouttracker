package v1

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(rg *gin.RouterGroup, handler *Handler) {
	user := rg.Group("/user")
	user.POST("/register", handler.Register)
	user.POST("/login", handler.Login)
}
