package models

type Level struct {
	ID    int    `json:"id" example:"1"`
	Value string `json:"value" example:"beginner"`
	Label string `json:"label" example:"Beginner"`
}
