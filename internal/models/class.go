package models

import (
	"github.com/google/uuid"
	"time"
)

type Class struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Key          string          `json:"key"`
	TeacherId    uuid.UUID       `json:"teacher_id"`
	Day          string          `json:"date"`
	StartTime    time.Time       `json:"start_time"`
	EndTime      time.Time       `json:"end_time"`
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
