package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/constants"
	"github.com/globus303/sportujspolu/models"
)

type userResponse struct {
	ID     int64  `json:"id" example:"123"`
	Name   string `json:"name" example:"John Doe"`
	Email  string `json:"email" example:"email@test.com"`
	Rating int    `json:"rating" example:"3"`
}

// @Summary Get current user
// @Description Gets the current user
// @Tags user
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} userResponse
// @Router /user/me [get]
func (s *Service) GetMe(c *gin.Context) {
	userID, _ := c.Get(constants.UserID)

	u := models.User{}

	err := s.db.QueryRow(`SELECT ID, Name, Email, Rating FROM users WHERE ID = ?`, userID).Scan(&u.ID, &u.Name, &u.Email, &u.Rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse := userResponse{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Rating: u.Rating,
	}

	c.JSON(http.StatusOK, userResponse)
}
