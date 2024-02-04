package models

type SectionClass struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Material    []Material   `json:"material"`
	Submission  []Submission `json:"submission"`
	ClassId     int          `json:"class_id"`
	Order       int          `json:"order"`
}
