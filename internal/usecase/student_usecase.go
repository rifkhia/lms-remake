package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type StudentUsecase interface {
	FetchStudentById(c context.Context, id uuid.UUID) (*models.Student, error)
	FetchStudentByName(c context.Context, name string) ([]*models.Student, error)
	Register(c context.Context, student *models.LoginRegisterRequest) (interface{}, error)
	Login(c context.Context, request *models.LoginRegisterRequest) (interface{}, error)
}

type StudentUsecaseImpl struct {
	studentRepo repository.StudentRepository
}

func (s *StudentUsecaseImpl) FetchStudentById(c context.Context, id uuid.UUID) (*models.Student, error) {
	studentResult, err := s.studentRepo.GetStudentByID(c, id)
	if err != nil {
		return nil, err
	}

	return studentResult, nil
}

func (s *StudentUsecaseImpl) FetchStudentByName(c context.Context, name string) ([]*models.Student, error) {
	studentResult, err := s.studentRepo.GetStudentByName(c, name)
	if err != nil {
		return nil, err
	}

	if studentResult == nil {
		return nil, err
	}

	return studentResult, nil
}

func (s *StudentUsecaseImpl) Register(c context.Context, student *models.LoginRegisterRequest) (interface{}, error) {
	studentRequest, err := student.NewStudent()
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.CreateAccessToken(studentRequest.ID, utils.STUDENT_ROLE)
	if err != nil {
		return nil, errors.New("Failed to create access token")
	}

	refreshToken, err := utils.CreateRefreshToken(studentRequest.ID, accessToken, utils.STUDENT_ROLE)
	if err != nil {
		return nil, errors.New("Failed to refresh token")
	}

	err = s.studentRepo.CreateStudent(c, studentRequest)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Access token":  accessToken,
		"Refresh token": refreshToken,
	}, nil
}

func (s *StudentUsecaseImpl) Login(c context.Context, request *models.LoginRegisterRequest) (interface{}, error) {
	if request.Email == "" || request.Password == "" {
		return nil, errors.New("Email or password required")
	}

	student, err := s.studentRepo.GetStudentByEmail(c, request.Email)
	if err != nil {
		return nil, err
	}

	err = utils.ValidatePassword(student.Password, request.Password)
	if err != nil {
		return nil, errors.New("Password missmatch")
	}

	accessToken, err := utils.CreateAccessToken(student.ID, utils.STUDENT_ROLE)
	if err != nil {
		return nil, errors.New("Failed to create access token")
	}

	refreshToken, err := utils.CreateRefreshToken(student.ID, accessToken, utils.STUDENT_ROLE)
	if err != nil {
		return nil, errors.New("Failed to refresh token")
	}

	return map[string]interface{}{
		"Access token":  accessToken,
		"Refresh token": refreshToken,
	}, nil
}

func NewStudentUsecase(repo repository.StudentRepository) StudentUsecase {
	return &StudentUsecaseImpl{
		studentRepo: repo,
	}
}
