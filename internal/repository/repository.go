package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/models"
)

type StudentRepository interface {
	GetStudentByID(c context.Context, id uuid.UUID) (*models.Student, error)
	GetStudentByEmail(c context.Context, email string) (*models.Student, error)
	GetStudentByName(c context.Context, name string) ([]*models.Student, error)
	CreateStudent(c context.Context, student *models.Student) error
}
