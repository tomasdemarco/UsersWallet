package main

import (
	"PruebaGo/models"
	"PruebaGo/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := connectDB()
	if err != nil {
		return
	}

	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Use(dbMiddleware(*conn))

	usersGroup := router.Group("users")
	{
		usersGroup.POST("register", routes.UserRegister)
		usersGroup.POST("login", routes.UsersLogin)
		usersGroup.GET("allusers", authMiddleWare(), routes.GetAllUsers)
		usersGroup.GET("user", authMiddleWare(), routes.GetUser)
		usersGroup.DELETE("delete", routes.UserDelete)
		usersGroup.PUT("update", routes.UserUpdate)
	}

	router.Run(":" + os.Getenv("PORT"))
}

func connectDB() (c *pgx.Conn, err error) {

	dbConn := os.Getenv("DBCONNECTION")
	conn, err := pgx.Connect(context.Background(), dbConn)
	if err != nil || conn == nil {
		fmt.Println("Error connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Set("db", conn)
		c.Next()
	}
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
			return
		}
		token := split[1]
		//fmt.Printf("Bearer (%v) \n", token)
		isValid, userID := models.IsTokenValid(token)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}
