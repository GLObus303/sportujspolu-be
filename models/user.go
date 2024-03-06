package models

type User struct {
	ID int64 `json:"id,omitempty" example:"123"`
	// Public_ID   string    `json:"id" example:"pwnrxtbi9z0v"`
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"email@test.com"`
	Password string `json:"password,omitempty" example:"Test123"`
	Rating   int    `json:"rating" example:"3"`
}
