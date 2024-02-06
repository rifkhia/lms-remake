package models

type Submission struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	File           string `json:"file"`
	ClassSectionId int    `json:"class_section_id"`
}
