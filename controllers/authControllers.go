package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rwiteshbera/MoneyTracker/config"
	"github.com/rwiteshbera/MoneyTracker/models"
	"github.com/rwiteshbera/MoneyTracker/utils"
	"net/http"
)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database, err := config.ConnectDB()
		defer database.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (firstName TEXT, lastName TEXT, phoneNumber INTEGER PRIMARY KEY, password TEXT)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, err = statement.Exec()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		statement, err = database.Prepare("INSERT INTO users (firstName, lastName, phoneNumber, password) VALUES (?, ?, ?, ?)")
		defer statement.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Encrypt the password with bcryptjs
		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = statement.Exec(user.FirstName, user.LastName, user.PhoneNumber, user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "registration successful!"})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User // Storing user input in this instance

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database, err := config.ConnectDB()
		defer database.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		statement, err := database.Prepare("SELECT firstName, lastName, phoneNumber, password FROM users WHERE phoneNumber = ?")
		defer statement.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var savedUser models.User // Storing the saved user data in this instance
		err = statement.QueryRow(user.PhoneNumber).Scan(&savedUser.FirstName, &savedUser.LastName, &savedUser.PhoneNumber, &savedUser.Password)

		// If no user found
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No user found"})
			return
		}

		// If password valid
		isPasswordValid, err := utils.VerifyPassword(user.Password, savedUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !isPasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		token, err := utils.GenerateToken(savedUser.FirstName, savedUser.LastName, savedUser.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
