package routes

import (
	"PruebaGo/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func UsersLogin(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.IsAuthenticated(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "There was an error authenticating.",
	})
}

func UserRegister(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.RegisterUser(&conn)
	if err != nil {
		fmt.Println("Error in user.Register()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// token, err := user.GetAuthToken()
	// if err == nil {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"token": token,
	// 	})
	// 	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"Registrado": "Usuario generado con exito.",
	})
}
func GetAllUsers(c *gin.Context) {

	user := models.User{}
	var users []models.User
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	users = user.AllUsers(&conn)
	//c.JSON(http.StatusOK, gin.H{"recordset": &users})
	c.IndentedJSON(http.StatusOK, gin.H{"recordset": &users})
}

func GetUser(c *gin.Context) {

	user := models.User{}
	var users []models.User
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	users = user.UserGet(&conn, c)
	//c.JSON(http.StatusOK, gin.H{"recordset": &users})
	c.IndentedJSON(http.StatusOK, gin.H{"user": &users})
}

func UserDelete(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.DeleteUser(&conn)

	if err != nil {
		fmt.Println("Error in user.Delete()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Delete user_id": user.ID,
	})
}

func UserUpdate(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.UpdateUser(&conn)
	if err != nil {
		fmt.Println("Error in user.Update()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
	})
}
