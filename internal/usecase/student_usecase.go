package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/dto"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
	"log"
	"strings"
)

type StudentUsecase interface {
	FetchStudentById(c context.Context, id uuid.UUID) (*models.StudentProfile, pkg.CustomError)
	FetchStudentByName(c context.Context, name string) ([]*models.Student, pkg.CustomError)
	Register(c context.Context, student *dto.StudentRegisterRequest) (interface{}, pkg.CustomError)
	Login(c context.Context, request *dto.StudentLoginRequest) (interface{}, pkg.CustomError)
	DeleteStudent(c context.Context, id uuid.UUID) pkg.CustomError
	EditProfileStudent(c context.Context, request *dto.StudentProfileRequest) pkg.CustomError
	FetchStudentSchedule(c context.Context, id uuid.UUID) ([]*models.StudentSchedule, pkg.CustomError)
}

type StudentUsecaseImpl struct {
	studentRepo repository.StudentRepository
}

func (s *StudentUsecaseImpl) FetchStudentById(c context.Context, id uuid.UUID) (*models.StudentProfile, pkg.CustomError) {

	studentResult, err := s.studentRepo.GetStudentByID(c, id)
	if err.Cause != nil {
		return nil, err
	}
	studentProfile, err := s.studentRepo.GetStudentProfile(c, id)
	if err.Cause != nil {
		return nil, err
	}

	student := models.StudentProfile{
		ID:          studentResult.ID,
		Name:        studentResult.Name,
		NIM:         studentResult.NIM,
		Email:       studentResult.Email,
		DateOfBirth: studentProfile.DateOfBirth,
		Gender:      studentProfile.Gender,
		Address:     studentProfile.Address,
		Phone:       studentProfile.Phone,
	}

	return &student, pkg.CustomError{}
}

func (s *StudentUsecaseImpl) FetchStudentByName(c context.Context, name string) ([]*models.Student, pkg.CustomError) {
	studentResult, err := s.studentRepo.GetStudentByName(c, name)
	if err.Cause != nil {
		return nil, err
	}

	return studentResult, pkg.CustomError{}
}

func (s *StudentUsecaseImpl) Register(c context.Context, student *dto.StudentRegisterRequest) (interface{}, pkg.CustomError) {
	studentRequest, err := student.NewStudent()
	if err.Cause != nil {
		return nil, err
	}

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

func (s *StudentUsecaseImpl) Login(c context.Context, request *dto.StudentLoginRequest) (interface{}, pkg.CustomError) {
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

func (s *StudentUsecaseImpl) DeleteStudent(c context.Context, id uuid.UUID) pkg.CustomError {
	customError := s.studentRepo.DeleteStudent(c, id)
	if customError.Cause != nil {
		return customError
	}

	return pkg.CustomError{}
}

func (s *StudentUsecaseImpl) EditProfileStudent(c context.Context, request *dto.StudentProfileRequest) pkg.CustomError {
	student, customError := s.studentRepo.GetStudentProfile(c, request.ID)
	if customError.Cause != nil {
		if strings.Contains(customError.Cause.Error(), "no student profile found") {
			customError = s.studentRepo.AddStudentProfile(c, request)
			if customError.Cause != nil {
				return customError
			}
			return pkg.CustomError{}
		}
		return customError
	}

	request.UpdateStudent(student)
	log.Println(student)
	customError = s.studentRepo.EditStudentProfile(c, student)
	if customError.Cause != nil {
		return customError
	}

	return pkg.CustomError{}
}

func (s *StudentUsecaseImpl) FetchStudentSchedule(c context.Context, id uuid.UUID) ([]*models.StudentSchedule, pkg.CustomError) {
	studentSchedules, customError := s.studentRepo.FetchStudentClass(c, id)
	if customError.Cause != nil {
		return nil, customError
	}

	for _, schedule := range studentSchedules {
		schedule.Day = utils.ConvertIntToDay(schedule.Day)
		schedule.StartTime, customError = utils.ConvertTimes(schedule.StartTime)
		if customError.Cause != nil {
			return nil, customError
		}
		schedule.EndTime, customError = utils.ConvertTimes(schedule.EndTime)
		if customError.Cause != nil {
			return nil, customError
		}
	}

	return studentSchedules, pkg.CustomError{}
}

func NewStudentUsecase(repo repository.StudentRepository) StudentUsecase {
	return &StudentUsecaseImpl{
		studentRepo: repo,
	}
}
