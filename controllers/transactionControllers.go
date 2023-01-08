package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/MoneyTracker/config"
	"github.com/rwiteshbera/MoneyTracker/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction models.Transaction

		transaction.Id = uint64(time.Now().Unix())
		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		i, err := strconv.ParseUint(transaction.CreatedBy, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		transaction.Id = uint64(time.Now().Unix()) | i
		database, err1 := config.ConnectDB()
		defer database.Close()
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}
		statement, err2 := database.Prepare("CREATE TABLE IF NOT EXISTS transactions (tid INTEGER NOT NULL, transactionName TEXT NOT NULL, amount INTEGER NOT NULL,createdBy TEXT NOT NULL,FOREIGN KEY (createdBy) REFERENCES users(phoneNumber) ON DELETE CASCADE)")
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}
		_, err3 := statement.Exec()
		if err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
			return
		}

		statement, err4 := database.Prepare("INSERT INTO transactions (tid, transactionName, amount, createdBy) VALUES (?, ?, ?, ?)")
		defer statement.Close()

		if err4 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err4.Error()})
			return
		}

		_, err5 := statement.Exec(transaction.Id, transaction.TransactionName, transaction.Amount, transaction.CreatedBy)
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

		statement, err := database.Prepare("DELETE FROM transactions WHERE tid = ?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		_, err = statement.Exec(transaction.Id)
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

		rows, err := database.Query("SELECT tid, transactionName, amount, createdBy FROM transactions WHERE createdBy = ?", userPhoneNumber)
		defer rows.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var transactions []models.Transaction
		for rows.Next() {
			var transactionCreatedByU models.Transaction

			err := rows.Scan(&transactionCreatedByU.Id, &transactionCreatedByU.TransactionName, &transactionCreatedByU.Amount, &transactionCreatedByU.CreatedBy)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			transactions = append(transactions, transactionCreatedByU)
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
		userPhoneNumber, _ := c.Get("phoneNumber")
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

		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS members (phoneNumber INTEGER NOT NULL, firstname TEXT NOT NULL, lastname TEXT NOT NULL, amountToBePaid INTEGER NOT NULL,transactionId INTEGER NOT NULL,createdBy TEXT NOT NULL,FOREIGN KEY (phoneNumber) REFERENCES users(phoneNumber) ON DELETE CASCADE,FOREIGN KEY (transactionId) REFERENCES transactions(tid) ON DELETE CASCADE,FOREIGN KEY (createdBy) REFERENCES users(phoneNumber) ON DELETE CASCADE,PRIMARY KEY(phoneNumber, transactionId))")
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

		statement, err = database.Prepare("INSERT INTO members (phoneNumber, firstname, lastname, amountToBePaid, transactionId, createdBy) VALUES (?, ?, ?, ?, ?, ?)")
		defer statement.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = statement.Exec(newMember.PhoneNumber, newMember.FirstName, newMember.LastName, amount/(count+2), newMember.TransactionId, userPhoneNumber)
		if err != nil {
			// If member already exists
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Member already exists in this transaction."})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
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
		userPhoneNumber, _ := c.Get("phoneNumber")

		database, err := config.ConnectDB()
		defer database.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := database.Query("SELECT phoneNumber, firstname, lastname, amountToBePaid, transactionId FROM members WHERE createdBy = ?", userPhoneNumber)
		defer rows.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var members []models.Member
		for rows.Next() {
			var member models.Member

			err := rows.Scan(&member.PhoneNumber, &member.FirstName, &member.LastName, &member.AmountToBePaid, &member.TransactionId)
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

// Function Search member by phone number
func SearchMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database, err1 := config.ConnectDB()
		defer database.Close()
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}

		if len(user.PhoneNumber) != 10 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid phone number"})
			return
		}
		statement, err := database.Prepare("SELECT firstName, lastName, phoneNumber FROM users WHERE phoneNumber = ?")
		defer statement.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var savedUser models.User // Storing the saved user data in this instance
		err = statement.QueryRow(user.PhoneNumber).Scan(&savedUser.FirstName, &savedUser.LastName, &savedUser.PhoneNumber)

		// If no user found
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No user found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": savedUser})
	}
}
