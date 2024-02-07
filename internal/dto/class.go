package dto

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/utils"
	"time"
)

type ClassCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Day         string `json:"day"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type ClassByNameResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Day         string `json:"day"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

func (c *ClassCreate) NewClass(teacherId uuid.UUID) (*models.Class, pkg.CustomError) {
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
	return &models.Class{
		Name:        c.Name,
		Description: c.Description,
		Key:         classKey,
		TeacherId:   teacherId,
		Day:         utils.ConvertDaysToInt(c.Day),
		StartTime:   fmt.Sprintf("%s", parsedStartTime),
		EndTime:     fmt.Sprintf("%s", parsedEndTime),
	}, pkg.CustomError{}
}
