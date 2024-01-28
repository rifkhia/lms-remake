package models

import "github.com/google/uuid"

type Teacher struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	NPM      string    `json:"NPM"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}
