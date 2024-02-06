package models

import (
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/utils"
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

type ClassCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Day         string `json:"day"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

func (c *ClassCreate) NewClass(teacherId uuid.UUID) (*Class, pkg.CustomError) {
	classKey := utils.GenerateClassKey(8)

	parsedStartTime, err := time.Parse("15:04", c.StartTime)
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Code:    utils.BAD_REQUEST,
			Service: utils.MODEL_SERVICE,
		}
	}

	parsedEndTime, err := time.Parse("15:04", c.EndTime)
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Code:    utils.BAD_REQUEST,
			Service: utils.MODEL_SERVICE,
		}
	}
	return &Class{
		Name:        c.Name,
		Description: c.Description,
		Key:         classKey,
		TeacherId:   teacherId,
		Day:         c.Day,
		StartTime:   parsedStartTime,
		EndTime:     parsedEndTime,
	}, pkg.CustomError{}
}

//type StudentClass struct {
//	ID        int       `json:"ID"`
//	ClassId   int       `json:"class_id"`
//	StudentId uuid.UUID `json:"student_id"`
//}
