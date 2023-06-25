package user

import (
	"database/sql"
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *Service {
	return &Service{db}
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// func VerifyPassword(password, hashedPassword string) error {
// 	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// }

// func loginCheck(username string, password string, s *Service) (string, error) {

// 	var err error

// 	u := User{}

// 	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

// 	if err != nil {
// 		return "", err
// 	}

// 	err = VerifyPassword(password, u.Password)

// 	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
// 		return "", err
// 	}

// 	token, err := token.GenerateToken(u.ID)

// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil

// }

func (s *Service) Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}

	u.Username = input.Username
	u.Password = input.Password

	// token, err := LoginCheck(u.Username, u.Password)

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"token": token})
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Service) Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}

	u.Username = input.Username
	u.Password = input.Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
