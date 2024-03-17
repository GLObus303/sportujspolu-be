package models

type User struct {
	ID       string `json:"id" example:"pwnrxtbi9z0v"`
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"email@test.com"`
	Password string `json:"-" example:"Test123"`
	Rating   int    `json:"rating" example:"3"`
}

type PublicUser struct {
	ID     string `json:"id" example:"pwnrxtbi9z0v"`
	Name   string `json:"name" example:"John Doe"`
	Email  string `json:"email" example:"email@test.com"`
	Rating int    `json:"rating" example:"3"`
}
