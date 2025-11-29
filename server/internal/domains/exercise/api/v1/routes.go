package v1

import "github.com/gin-gonic/gin"

func RegisterExerciseRoutes(rg *gin.RouterGroup, handler *Handler, auth gin.HandlerFunc) {
	user := rg.Group("/exercise")
	user.POST("/create", auth, handler.CreateExercise)
	user.GET("/getForUser", auth, handler.GetExerciseForUser)
}
