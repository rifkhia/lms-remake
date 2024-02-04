package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type StudentUsecase interface {
	FetchStudentById(c context.Context, id uuid.UUID) (*models.Student, pkg.CustomError)
	FetchStudentByName(c context.Context, name string) ([]*models.Student, pkg.CustomError)
	Register(c context.Context, student *models.StudentRegisterRequest) (interface{}, pkg.CustomError)
	Login(c context.Context, request *models.StudentLoginRequest) (interface{}, pkg.CustomError)
}

type StudentUsecaseImpl struct {
	studentRepo repository.StudentRepository
}

func (s *StudentUsecaseImpl) FetchStudentById(c context.Context, id uuid.UUID) (*models.Student, pkg.CustomError) {

	studentResult, err := s.studentRepo.GetStudentByID(c, id)
	if err.Cause != nil {
		return nil, err
	}

	return studentResult, pkg.CustomError{}
}

func (s *StudentUsecaseImpl) FetchStudentByName(c context.Context, name string) ([]*models.Student, pkg.CustomError) {
	studentResult, err := s.studentRepo.GetStudentByName(c, name)
	if err.Cause != nil {
		return nil, err
	}

	return studentResult, pkg.CustomError{}
}

func (s *StudentUsecaseImpl) Register(c context.Context, student *models.StudentRegisterRequest) (interface{}, pkg.CustomError) {
	studentRequest, err := student.NewStudent()
	if err.Cause != nil {
		return nil, err
	}

	//studentCheck, err := s.studentRepo.GetStudentByEmail(c, student.Email)
	//if studentCheck == nil {
	//	customError := pkg.CustomError{
	//		Cause:   errors.New("email already exists"),
	//		Code:    utils.BAD_REQUEST,
	//		Service: utils.USECASE_SERVICE,
	//	}
	//	return nil, customError
	//}

	err = s.studentRepo.CreateStudent(c, studentRequest)
	if err.Cause != nil {
		return nil, err
	}

	accessToken, err := utils.CreateAccessToken(studentRequest.ID, utils.STUDENT_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	refreshToken, err := utils.CreateRefreshToken(studentRequest.ID, accessToken, utils.STUDENT_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Access token":  accessToken,
		"Refresh token": refreshToken,
	}, pkg.CustomError{}
}

func (s *StudentUsecaseImpl) Login(c context.Context, request *models.StudentLoginRequest) (interface{}, pkg.CustomError) {
	if request.Email == "" || request.Password == "" {
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("email or password can't be blank"),
			Service: utils.USECASE_SERVICE,
		}
	}

	student, err := s.studentRepo.GetStudentByEmail(c, request.Email)
	if err.Cause != nil {
		return nil, err
	}

	err = utils.ValidatePassword(student.Password, request.Password)
	if err.Cause != nil {
		return nil, err
	}

	accessToken, err := utils.CreateAccessToken(student.ID, utils.STUDENT_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	refreshToken, err := utils.CreateRefreshToken(student.ID, accessToken, utils.STUDENT_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Access token":  accessToken,
		"Refresh token": refreshToken,
	}, pkg.CustomError{}
}

func NewStudentUsecase(repo repository.StudentRepository) StudentUsecase {
	return &StudentUsecaseImpl{
		studentRepo: repo,
	}
}
