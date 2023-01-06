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
		transaction.Id = (uint64(time.Now().Unix()) | transaction.CreatedBy)

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
		statement, err2 := database.Prepare("CREATE TABLE IF NOT EXISTS transactions (tid INTEGER NOT NULL,amount INTEGER NOT NULL,createdBy INTEGER NOT NULL,FOREIGN KEY (createdBy) REFERENCES users(phoneNumber) ON DELETE CASCADE)")
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}
		_, err3 := statement.Exec()
		if err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
			return
		}

		statement, err4 := database.Prepare("INSERT INTO transactions (tid, amount, createdBy) VALUES (?, ?, ?)")
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

// Delete A transaction
func DeleteTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionId uint64

		database, err := config.ConnectDB()
		defer database.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		statement, err := database.Prepare("DELETE FROM transactions WHERE tid = ?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		_, err = statement.Exec(transactionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Transaction has been deleted!"})
	}
}

// See transactions created by a loggedIn user
func ShowTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userPhoneNumber, _ := c.Get("phoneNumber")

		database, err := config.ConnectDB()
		defer database.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := database.Query("SELECT tid, amount, createdBy FROM transactions WHERE createdBy = ?", userPhoneNumber)
		defer rows.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var transactions []models.Transaction
		for rows.Next() {
			var transaction models.Transaction

			err := rows.Scan(&transaction.Id, &transaction.Amount, &transaction.CreatedBy)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			transactions = append(transactions, transaction)
		}

		err = rows.Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": transactions})
	}
}

// Add member in a transaction
func AddMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newMember models.Member

		if err := c.BindJSON(&newMember); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		database, err1 := config.ConnectDB()
		defer database.Close()
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}

		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS members (phoneNumber INTEGER NOT NULL,amountToBePaid INTEGER NOT NULL, transactionId INTEGER NOT NULL,FOREIGN KEY (phoneNumber) REFERENCES users(phoneNumber) ON DELETE CASCADE,FOREIGN KEY (transactionId) REFERENCES transactions(tid) ON DELETE CASCADE)")
		defer statement.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, err = statement.Exec()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var count uint64
		var amount uint64
		err = database.QueryRow("SELECT COUNT(*) FROM members WHERE transactionId=?", newMember.TransactionId).Scan(&count)
		err = database.QueryRow("SELECT amount FROM transactions WHERE tid=?", newMember.TransactionId).Scan(&amount)

		statement, err = database.Prepare("INSERT INTO members (phoneNumber, amountToBePaid, transactionId) VALUES (?, ?, ?)")
		defer statement.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = statement.Exec(newMember.PhoneNumber, amount/(count+2), newMember.TransactionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//UPDATE members SET amountToBePaid = ?
		statement, err = database.Prepare("UPDATE members SET amountToBePaid = ? WHERE transactionId = ?")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = statement.Exec(amount/(count+2), newMember.TransactionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"count": count, "message": "member has been added"})
	}
}

func ShowMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction models.Transaction

		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		database, err := config.ConnectDB()
		defer database.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := database.Query("SELECT phoneNumber, amountToBePaid, transactionId FROM members WHERE transactionId = ?", &transaction.Id)
		defer rows.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var members []models.Member
		for rows.Next() {
			var member models.Member

			err := rows.Scan(&member.PhoneNumber, &member.AmountToBePaid, &member.TransactionId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			members = append(members, member)
		}

		err = rows.Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": members})
	}
}

func DeleteMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newMember models.Member

		if err := c.BindJSON(&newMember); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		database, err1 := config.ConnectDB()
		defer database.Close()
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}

		statement, err := database.Prepare("DELETE FROM members WHERE phoneNumber = ? AND transactionId = ?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		_, err = statement.Exec(newMember.PhoneNumber, newMember.TransactionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Member has been deleted!"})
	}
}

func MarkAsPaid() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
