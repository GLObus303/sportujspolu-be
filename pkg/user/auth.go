package user

import (
	"database/sql"
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/models"
	"github.com/globus303/sportujspolu/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *Service {
	return &Service{db}
}

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func loginCheck(email string, password string, s *Service) (string, error) {
	var err error

	u := models.User{}

	err = s.db.QueryRow(`SELECT * FROM users WHERE email = $1`, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Rating)
	if err != nil {
		return "", err
	}

	err = verifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := utils.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}

type LoginInput struct {
	Email    string `json:"email" binding:"required" example:"email@test.com"`
	Password string `json:"password" binding:"required" example:"Test123"`
}

// @Summary User login
// @Description Logs in a user with the provided credentials
// @Tags user
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login credentials"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /user/login [post]
func (s *Service) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Email = input.Email
	u.Password = input.Password

	token, err := loginCheck(u.Email, u.Password, s)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required" example:"email@test.com"`
	Password string `json:"password" binding:"required" example:"Test123"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
}

// @Summary Register a new user
// @Description Registers a new user with the provided details
// @Tags user
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Registration details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /user/register [post]
func (s *Service) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Name = input.Name
	u.Email = input.Email
	u.Password = input.Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.Password = string(hashedPassword)
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	query := `INSERT INTO users (name, email, password, rating) VALUES ($1, $2, $3, 0)`
	_, err = s.db.Exec(query, u.Name, u.Email, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
