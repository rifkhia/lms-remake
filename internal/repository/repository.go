package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/dto"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
)

type StudentRepository interface {
	GetStudentByID(c context.Context, id uuid.UUID) (*models.Student, pkg.CustomError)
	GetStudentByEmail(c context.Context, email string) (*models.Student, pkg.CustomError)
	GetStudentByName(c context.Context, name string) ([]*models.Student, pkg.CustomError)
	CreateStudent(c context.Context, student *models.Student) pkg.CustomError
	DeleteStudent(c context.Context, id uuid.UUID) pkg.CustomError
	GetStudentByClassId(c context.Context, classId int) ([]*models.StudentClass, pkg.CustomError)
	GetStudentProfile(c context.Context, id uuid.UUID) (*dto.StudentProfileRequest, pkg.CustomError)
	AddStudentProfile(c context.Context, student *dto.StudentProfileRequest) pkg.CustomError
	EditStudentProfile(c context.Context, student *dto.StudentProfileRequest) pkg.CustomError
	FetchStudentClass(c context.Context, id uuid.UUID) ([]*models.StudentSchedule, pkg.CustomError)
}

type ClassRepository interface {
	GetClassByID(c context.Context, id int) (*models.Class, pkg.CustomError)
	GetClassByTeacherID(c context.Context, teacherId int) ([]*models.Class, pkg.CustomError)
	GetClassByName(c context.Context, name string) ([]*models.Class, pkg.CustomError)
	CreateClass(c context.Context, class *models.Class) pkg.CustomError
	DeleteClass(c context.Context, id int) pkg.CustomError
	JoinClass(c context.Context, classId int, studentId uuid.UUID) pkg.CustomError
	LeftCLass(c context.Context, classId int, studentId uuid.UUID) pkg.CustomError
	CheckStudentClassExists(c context.Context, classId int, studentId uuid.UUID) (bool, pkg.CustomError)
	CheckTeacherClassExists(c context.Context, teacherId uuid.UUID, classId int) (bool, pkg.CustomError)
	GetClassSectionByClassId(c context.Context, classId int) ([]*models.SectionClass, pkg.CustomError)
	CreateClassSection(c context.Context, classSection *models.SectionClass) pkg.CustomError
	InsertSubmissionTeacher(c context.Context, request *models.Submission) pkg.CustomError
	GetClassSectionById(c context.Context, id int) (*models.SectionClass, pkg.CustomError)
	UpdateClassSection(c context.Context, classSection *models.SectionClass) pkg.CustomError
	InsertSubmissionStudent(c context.Context, request *dto.StudentSubmissionRequest) pkg.CustomError
	GetSubmissionByClassSection(c context.Context, classSecctionId int) ([]*models.StudentSubmission, pkg.CustomError)
}

type TeacherRepository interface {
	GetTeacherByEmail(c context.Context, email string) (*models.Teacher, pkg.CustomError)
	GetTeacherById(c context.Context, id uuid.UUID) (*models.Teacher, pkg.CustomError)
	CreateTeacher(c context.Context, request *models.Teacher) pkg.CustomError
}
