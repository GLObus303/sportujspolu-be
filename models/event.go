package models

import "time"

type Event struct {
	ID          int64     `json:"-"`
	Public_ID   string    `json:"id" example:"pwnrxtbi9z0v"`
	Name        string    `json:"name" example:"Basketball Match at Park"`
	Sport       string    `json:"sport" example:"Basketball"`
	Date        time.Time `json:"date" example:"2023-07-10"`
	Location    string    `json:"location" example:"Central Park"`
	Price       float64   `json:"price" example:"0.00"`
	Description string    `json:"description" example:"Example Description"`
	Level       string    `json:"level" example:"Any"`
}
