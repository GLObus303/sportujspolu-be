package models

import "time"

type Event struct {
	ID          uint16    `json:"-"`
	Public_ID   string    `json:"id" example:"pwnrxtbi9z0v"`
	Name        string    `json:"name" example:"Basketball Match at Park"`
	Sport       string    `json:"sport" example:"Basketball"`
	Date        time.Time `json:"date" example:"2023-11-03T10:15:30Z"`
	Location    string    `json:"location" example:"Central Park"`
	Price       uint16    `json:"price" example:"123"`
	Description string    `json:"description" example:"Example Description"`
	Level       string    `json:"level" example:"Any"`
	Created_At  time.Time `json:"createdAt" example:"2023-11-03T10:15:30Z"`
	Owner_ID    uint16    `json:"ownerId" example:"123"`
}
