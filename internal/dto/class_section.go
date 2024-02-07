package dto

import "github.com/rifkhia/lms-remake/internal/models"

type ClassSectionUpdate struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Order       int          `json:"order"`
	Task        bool         `json:"task"`
}

func (c *ClassSectionUpdate) UpdateClassSection(section *models.SectionClass)  {
	if c.Task != false {
		section.Task = true
	}

	if c.Title != "" {
		section.Title = c.Title
	}

	if c.Description != "" {
		section.Description = c.Description
	}

	if c.Order != 0 {
		section.Order = c.Order
	}
}