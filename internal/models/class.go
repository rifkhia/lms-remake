package models

import (
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type Class struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Key          string          `json:"key"`
	TeacherId    uuid.UUID       `json:"teacher_id"`
	ClassSection []*SectionClass `json:"class_section"`
	Student      []*StudentClass `json:"student"`
}

type ClassCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *ClassCreate) NewClass(teacherId uuid.UUID) Class {
	classKey := utils.GenerateClassKey(8)

	return Class{
		Name:        c.Name,
		Description: c.Description,
		Key:         classKey,
		TeacherId:   teacherId,
	}
}

//type StudentClass struct {
//	ID        int       `json:"ID"`
//	ClassId   int       `json:"class_id"`
//	StudentId uuid.UUID `json:"student_id"`
//}
