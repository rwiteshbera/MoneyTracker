package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/MoneyTracker/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controllers.Signup())
	incomingRoutes.POST("/login", controllers.Login())
}
