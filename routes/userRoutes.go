package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/MoneyTracker/controllers"
	"github.com/rwiteshbera/MoneyTracker/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middlewares.Authenticate())
	incomingRoutes.POST("/create", controllers.CreateTransaction())   // Create a new transaction
	incomingRoutes.POST("/delete", controllers.DeleteTransaction())   // Delete a new transaction
	incomingRoutes.GET("/get", controllers.ShowTransactions())        // View your transactions
	incomingRoutes.POST("/add_member", controllers.AddMember())       //Add member in a transaction
	incomingRoutes.GET("/show", controllers.ShowMembers())            // Show members of a transaction
	incomingRoutes.POST("/delete_member", controllers.DeleteMember()) // Delete a member from a transaction
	incomingRoutes.POST("/mark", controllers.MarkAsPaid())            // Mark a member as paid
	incomingRoutes.POST("/search", controllers.SearchMember())        // Search member by phone number
}
