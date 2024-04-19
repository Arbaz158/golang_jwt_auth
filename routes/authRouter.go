package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt-project/controllers"
)

func AuthRoutes(incomningRoutes *gin.Engine) {
	incomningRoutes.POST("/users/signup", controllers.Singup())
	incomningRoutes.POST("/users/login", controllers.Login())
}
