package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/constants"
	"github.com/globus303/sportujspolu/models"
)

// @Summary Get current user
// @Description Gets the current user
// @Tags user
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} models.PublicUser
// @Router /user/me [get]
func (s *Service) GetMe(c *gin.Context) {
	userID, _ := c.Get(constants.UserID_key)

	u := models.User{}

	err := s.db.QueryRow(`SELECT id, name, email, rating FROM users WHERE ID = $1`, userID).Scan(&u.ID, &u.Name, &u.Email, &u.Rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse := models.PublicUser{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Rating: u.Rating,
	}

	c.JSON(http.StatusOK, userResponse)
}

// @Summary Delete current user
// @Description Deletes the current user
// @Tags user
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} string
// @Router /user/me [delete]
func (s *Service) DeleteMe(c *gin.Context) {
	userID, _ := c.Get(constants.UserID_key)

	_, err := s.db.Exec(`DELETE FROM users WHERE ID = $1`, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = s.db.Exec(`DELETE FROM events WHERE owner_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
