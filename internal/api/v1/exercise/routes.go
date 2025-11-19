package exercise

import "github.com/gin-gonic/gin"

func RegisterExerciseRoutes(rg *gin.RouterGroup, handler *Handler, auth gin.HandlerFunc) {
	user := rg.Group("/exercise")
	user.POST("/create", auth, handler.CreateExercise)
	user.POST("/getForUser", auth, handler.GetExerciseForUser)
}
