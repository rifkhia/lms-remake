package models

type Class struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}
