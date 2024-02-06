package models

import (
	"github.com/google/uuid"
	"time"
)

type Student struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	NIM      int       `json:"NIM"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type StudentProfile struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	NIM         int       `json:"NIM"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
}

type StudentClass struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type StudentSchedule struct {
	Name      string `json:"class_name"`
	Day       string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
