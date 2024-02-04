package models

import "github.com/google/uuid"

type Class struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Key          string          `json:"key"`
	TeacherId    string          `json:"teacher_id"`
	ClassSection []*SectionClass `json:"class_section"`
}

type StudentClass struct {
	ID        int       `json:"ID"`
	ClassId   int       `json:"class_id"`
	StudentId uuid.UUID `json:"student_id"`
}
