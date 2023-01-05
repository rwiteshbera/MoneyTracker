package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/MoneyTracker/controllers"
	"github.com/rwiteshbera/MoneyTracker/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middlewares.Authenticate())
	incomingRoutes.POST("/create", controllers.CreateTransaction())
}
