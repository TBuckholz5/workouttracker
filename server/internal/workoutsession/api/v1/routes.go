package v1

import "github.com/gin-gonic/gin"

func RegisterWorkoutSessionRoutes(rg *gin.RouterGroup, handler *Handler, auth gin.HandlerFunc) {
	route := rg.Group("/workoutsession")
	route.POST("/create", auth, handler.Create)
}
