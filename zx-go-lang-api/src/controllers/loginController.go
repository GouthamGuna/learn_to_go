package controllers

import (
	"fmt"
	"log"
	"net/http"
	"github.com/GouthamGuna/zx/src/services"
	"github.com/gin-gonic/gin"
)

func Login(r *gin.Engine) {

	r.POST("/auth/login", func(c *gin.Context) {

		username := c.PostForm("Username")
		password := c.PostForm("Password")

		fmt.Println("username : ", username)
		fmt.Println("password : ", password)

		if username == "" || password == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing credentials"})
            c.Abort()
            return
        }

		if username == "admin@zx" && password == "ZX#Auth@25" {
			tokenString, err := services.CreateToken(username)
			if err != nil {
				log.Println("Error creating token:", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Error to creating token"})
				c.Abort()
				return
			}
			log.Println("Token String:", tokenString)
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})
}
