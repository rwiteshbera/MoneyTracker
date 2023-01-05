package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/MoneyTracker/config"
	"github.com/rwiteshbera/MoneyTracker/models"
	"net/http"
	"time"
)

func CreateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		createdBy, _ := c.Get("phoneNumber")
		var transaction models.Transaction

		i, ok := createdBy.(uint64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "something went wrong"})
			return
		}
		transaction.CreatedBy = i
		transaction.Id = (uint64(time.Now().Unix() << 32)) | transaction.CreatedBy

		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database, err1 := config.ConnectDB()
		defer database.Close()
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}
		statement, err2 := database.Prepare("CREATE TABLE IF NOT EXISTS transactions (id INTEGER NOT NULL PRIMARY KEY,amount INTEGER NOT NULL,createdBy INTEGER NOT NULL,FOREIGN KEY (createdBy) REFERENCES users(phoneNumber) ON DELETE CASCADE)")
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}
		_, err3 := statement.Exec()
		if err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
			return
		}

		statement, err4 := database.Prepare("INSERT INTO transactions (id, amount, createdBy) VALUES (?, ?, ?)")
		defer statement.Close()

		if err4 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err4.Error()})
			return
		}

		_, err5 := statement.Exec(transaction.Id, transaction.Amount, transaction.CreatedBy)
		if err5 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err5.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
	}
}
