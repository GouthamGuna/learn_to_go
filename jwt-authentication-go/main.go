package main

import (
	jwt "github.com/GouthamGuna/jwt-authentication-go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Todo struct {
	Text string
	Done bool
}

var todos []Todo
var loggedInUser string

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Todos":    todos,
			"LoggedIn": loggedInUser != "",
			"Username": loggedInUser,
			"Role":     jwt.GetRole(loggedInUser),
		})
	})

	router.POST("/add", jwt.AuthenticateMiddleware, func(c *gin.Context) {
		text := c.PostForm("todo")
		todo := Todo{Text: text, Done: false}
		todos = append(todos, todo)
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/toggle", jwt.AuthenticateMiddleware, func(c *gin.Context) {
		index := c.PostForm("index")
		toggleIndex(index)
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Dummy credential check
		if (username == "employee" && password == "password") || (username == "senior" && password == "password") {
			tokenString, err := jwt.CreateToken(username)
			if err != nil {
				log.Println("Error creating token:", err)
				c.String(http.StatusInternalServerError, "Error creating token")
				return
			}

			log.Println("Token String:", tokenString)

			loggedInUser = username
			c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusSeeOther, "/")

		} else {
			c.String(http.StatusUnauthorized, "Invalid credentials")
		}
	})

	router.GET("/logout", func(c *gin.Context) {
		loggedInUser = ""
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/")																			
	})

	router.Run(os.Getenv("SERVER_PORT"))
}

func toggleIndex(index string) {
	i, _ := strconv.Atoi(index)
	if i >= 0 && i < len(todos) {
		todos[i].Done = !todos[i].Done
	}
}
