package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Lastname     string    `json:"lastname"`
	Document     int       `json:"document"`
	Birthday     string    `json:"birthday"`
	Email        string    `json:"email"`
	Phone        int       `json:"phone"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"-"`
}

var (
	tokenSecret = []byte(os.Getenv("KEYPRIVATE"))
)

func (u *User) RegisterUser(conn *pgx.Conn) error {
	u.Email = string(u.Email)

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("There was an error creating your account.")
	}
	u.PasswordHash = string(pwdHash)

	_, err = conn.Exec(context.Background(), "INSERT INTO users (name, lastname, documento, birthday, email, phone, password) VALUES($1, $2, $3, $4, $5, $6, $7)", u.Name, u.Lastname, u.Document, u.Birthday, u.Email, u.Phone, u.PasswordHash)
	return err
}

func (u *User) AllUsers(conn *pgx.Conn) []User {

	rows, err := conn.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rowSlice []User
	for rows.Next() {
		r := User{}
		err := rows.Scan(&r.ID, &r.Name, &r.Lastname, &r.Document, &r.Birthday, &r.Email, &r.Phone, &r.Password)
		if err != nil {
			log.Fatal(err)
		}
		rowSlice = append(rowSlice, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rowSlice
}

func (u *User) UserGet(conn *pgx.Conn, c *gin.Context) []User {

	userId, _ := c.Get("user_id")
	rows, err := conn.Query(context.Background(), "SELECT id, name, lastname, documento, birthday, email, phone FROM users where id = $1", userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rowSlice []User
	for rows.Next() {
		r := User{}
		err := rows.Scan(&r.ID, &r.Name, &r.Lastname, &r.Document, &r.Birthday, &r.Email, &r.Phone)
		if err != nil {
			log.Fatal(err)
		}
		rowSlice = append(rowSlice, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rowSlice
}

func (u *User) DeleteUser(conn *pgx.Conn) error {
	u.Email = string(u.Email)
	_, err := conn.Exec(context.Background(), "DELETE FROM users WHERE email=$1", u.Email)
	return err
}

func (u *User) UpdateUser(conn *pgx.Conn) error {
	u.Email = string(u.Email)
	u.Password = string(u.Password)
	_, err := conn.Exec(context.Background(), "UPDATE users SET email= $1, password=$2 WHERE email=$1", u.Email, u.Password)
	return err
}

// GetAuthToken returns the auth token to be used
func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
}

// IsAuthenticated checks to make sure password is correct and user is active
func (u *User) IsAuthenticated(conn *pgx.Conn) error {
	row := conn.QueryRow(context.Background(), "SELECT id, password from users WHERE email = $1", u.Email)
	err := row.Scan(&u.ID, &u.PasswordHash)
	if err == pgx.ErrNoRows {
		fmt.Println("User with email not found")
		return fmt.Errorf("Invalid login credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return fmt.Errorf("Invalid login credentials")
	}

	return nil
}

func IsTokenValid(tokenString string) (bool, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}

		return tokenSecret, nil
	})

	if err != nil {
		fmt.Printf("Err %v \n", err)
		return false, ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Println(claims)
		userID := claims["user_id"]
		return true, userID.(string)
	} else {
		fmt.Printf("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false, "uuid.UUID{}"
	}
}
