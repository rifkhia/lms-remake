package models

import (
	"github.com/google/uuid"
)

type Class struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Key          string          `json:"key"`
	TeacherId    uuid.UUID       `json:"teacher_id"`
	Day          string          `json:"day"`
	StartTime    string          `json:"start_time"`
	EndTime      string          `json:"end_time"`
	ClassSection []*SectionClass `json:"class_section"`
	Student      []*StudentClass `json:"student"`
}

type ClassWithIntDays struct {
}

//type StudentClass struct {
//	ID        int       `json:"ID"`
//	ClassId   int       `json:"class_id"`
//	StudentId uuid.UUID `json:"student_id"`
//}
