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

type ClassRepository interface {
	GetClassByID(c context.Context, id int) (*models.Class, error)
	GetClassByTeacherID(c context.Context, teacherId int) ([]*models.Class, error)
	GetClassByName(c context.Context, name string) ([]*models.Class, error)
	CreateClass(c context.Context, class *models.Class) error
	DeleteClass(c context.Context, id int) error
	JoinClass(c context.Context, classId int, studentId uuid.UUID) error
	CheckStudentClassExists(c context.Context, classId int, studentId uuid.UUID) (bool, error)
}
