package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt-project/controllers"
	"github.com/golang-jwt-project/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middlewares.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}
