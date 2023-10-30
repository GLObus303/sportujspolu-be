package models

type User struct {
	ID       int64  `json:"id" example:"123"`
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"email@test.com"`
	Password string `json:"password" example:"Test123"`
	Rating   int    `json:"rating" example:"3"`
}
