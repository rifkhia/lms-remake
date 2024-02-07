package models

import "time"

type Submission struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	File           string    `json:"file"`
	Deadline       time.Time `json:"deadline"`
	ClassSectionId int       `json:"class_section_id"`
}
