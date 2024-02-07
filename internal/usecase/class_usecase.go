package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/dto"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
	"github.com/spf13/viper"
	storage_go "github.com/supabase-community/storage-go"
	"mime/multipart"
)

type ClassUsecase interface {
	FetchClassById(c context.Context, id int) (*models.Class, pkg.CustomError)
	FetchClassByTeacherId(c context.Context, teacherId int) ([]*models.Class, pkg.CustomError)
	FetchClassByName(c context.Context, name string) ([]*dto.ClassByNameResponse, pkg.CustomError)
	CreateClass(c context.Context, request *dto.ClassCreate, teacherId uuid.UUID) pkg.CustomError
	CheckIfStudentInClass(c context.Context, studentId uuid.UUID, classId int) (bool, pkg.CustomError)
	JoinClass(c context.Context, studentId uuid.UUID, classId int, key string) pkg.CustomError
	LeftClass(c context.Context, classId int, studentId uuid.UUID) pkg.CustomError
	CreateSectionClass(c context.Context, request *models.SectionClass) pkg.CustomError
	AddSubmissionTeacher(c context.Context, request *models.Submission) pkg.CustomError
	AddSubmissionStudent(c context.Context, request *dto.StudentSubmissionRequest, file *multipart.FileHeader) pkg.CustomError
	FetchSectionClassById(c context.Context, id int) (*models.SectionClass, pkg.CustomError)
	FetchSubmissionBySection(c context.Context, sectionClassId int, teacherId uuid.UUID) ([]*models.StudentSubmission, pkg.CustomError)
}

type classUsecaseImpl struct {
	classRepo   repository.ClassRepository
	studentRepo repository.StudentRepository
}

func (s *classUsecaseImpl) FetchClassById(c context.Context, id int) (*models.Class, pkg.CustomError) {
	classResult, err := s.classRepo.GetClassByID(c, id)
	if err.Cause != nil {
		return nil, err
	}

	classResult.Day = utils.ConvertIntToDay(classResult.Day)
	classResult.StartTime, err = utils.ConvertTimes(classResult.StartTime)
	if err.Cause != nil {
		return nil, err
	}
	classResult.EndTime, err = utils.ConvertTimes(classResult.EndTime)
	if err.Cause != nil {
		return nil, err
	}

	classSection, err := s.classRepo.GetClassSectionByClassId(c, id)
	if err.Cause != nil {
		return nil, err
	}

	students, err := s.studentRepo.GetStudentByClassId(c, id)
	if err.Cause != nil {
		return nil, err
	}

	classResult.ClassSection = classSection
	classResult.Student = students

	return classResult, pkg.CustomError{}
}

func (s *classUsecaseImpl) FetchClassByTeacherId(c context.Context, teacherId int) ([]*models.Class, pkg.CustomError) {
	classResult, err := s.classRepo.GetClassByTeacherID(c, teacherId)
	if err.Cause != nil {
		return nil, err
	}

	return classResult, pkg.CustomError{}
}

func (s *classUsecaseImpl) FetchClassByName(c context.Context, name string) ([]*dto.ClassByNameResponse, pkg.CustomError) {
	var classes []*dto.ClassByNameResponse
	classResult, err := s.classRepo.GetClassByName(c, name)
	if err.Cause != nil {
		return nil, err
	}

	for _, c := range classResult {
		c.Day = utils.ConvertIntToDay(c.Day)
		c.StartTime, err = utils.ConvertTimes(c.StartTime)
		if err.Cause != nil {
			return nil, err
		}
		c.EndTime, err = utils.ConvertTimes(c.EndTime)
		if err.Cause != nil {
			return nil, err
		}
		class := dto.ClassByNameResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Day:         c.Day,
			StartTime:   c.StartTime,
			EndTime:     c.EndTime,
		}
		classes = append(classes, &class)
	}

	return classes, pkg.CustomError{}
}

func (s *classUsecaseImpl) CreateClass(c context.Context, request *dto.ClassCreate, teacherId uuid.UUID) pkg.CustomError {
	class, err := request.NewClass(teacherId)
	if err.Cause != nil {
		return err
	}

	err = s.classRepo.CreateClass(c, class)
	if err.Cause != nil {
		return err
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) CheckIfStudentInClass(c context.Context, studentId uuid.UUID, classId int) (bool, pkg.CustomError) {
	exists, err := s.classRepo.CheckStudentClassExists(c, classId, studentId)
	if err.Cause != nil {
		return false, err
	}

	if exists {
		return true, pkg.CustomError{}
	}

	return false, pkg.CustomError{}
}

func (s *classUsecaseImpl) JoinClass(c context.Context, studentId uuid.UUID, classId int, key string) pkg.CustomError {
	exists, err := s.classRepo.CheckStudentClassExists(c, classId, studentId)
	if err.Cause != nil {
		return err
	}

	if exists {
		return pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("you already join this class"),
			Service: utils.USECASE_SERVICE,
		}
	}

	classResult, err := s.classRepo.GetClassByID(c, classId)
	if err.Cause != nil {
		return err
	}

	if classResult.Key != key {
		return pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("missmatch class key"),
			Service: utils.USECASE_SERVICE,
		}
	}

	err = s.classRepo.JoinClass(c, classId, studentId)
	if err.Cause != nil {
		return err
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) FetchSectionClassById(c context.Context, id int) (*models.SectionClass, pkg.CustomError) {
	sectionClass, customError := s.classRepo.GetClassSectionById(c, id)
	if customError.Cause != nil {
		return nil, customError
	}

	return sectionClass, pkg.CustomError{}
}

func (s *classUsecaseImpl) CreateSectionClass(c context.Context, request *models.SectionClass) pkg.CustomError {
	err := s.classRepo.CreateClassSection(c, request)
	if err.Cause != nil {
		return err
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) LeftClass(c context.Context, classId int, studentId uuid.UUID) pkg.CustomError {
	customError := s.classRepo.LeftCLass(c, classId, studentId)
	if customError.Cause != nil {
		return customError
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) AddSubmissionTeacher(c context.Context, request *models.Submission) pkg.CustomError {
	customError := s.classRepo.InsertSubmissionTeacher(c, request)
	if customError.Cause != nil {
		return customError
	}

	sectionClass, customError := s.classRepo.GetClassSectionById(c, request.ClassSectionId)
	if customError.Cause != nil {
		return customError
	}

	sectionClass.Task = true

	customError = s.classRepo.UpdateClassSection(c, sectionClass)
	if customError.Cause != nil {
		return customError
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) AddSubmissionStudent(c context.Context, request *dto.StudentSubmissionRequest, file *multipart.FileHeader) pkg.CustomError {
	sectionClass, customError := s.classRepo.GetClassSectionById(c, request.ClassSectionId)
	if customError.Cause != nil {
		return customError
	}

	student, customError := s.studentRepo.GetStudentByID(c, request.ID)
	if customError.Cause != nil {
		return customError
	}

	if sectionClass.Task != true {
		return pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("there's currently no task here"),
			Service: utils.USECASE_SERVICE,
		}
	}

	parsedFile, err := file.Open()
	if err != nil {
		return pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.USECASE_SERVICE,
		}
	}

	storageGo := storage_go.NewClient(viper.GetString("SUPABASE_URL"), viper.GetString("SUPABASE_TOKEN"), nil)
	_, err = storageGo.UploadFile("submissions_student", fmt.Sprintf("class_section/%s/%s-%d.pdf", sectionClass.Title, student.Name, student.NIM), parsedFile)
	linkFile := storageGo.GetPublicUrl("submissions_student", fmt.Sprintf("class_section/%s/%s-%d.pdf", sectionClass.Title, student.Name, student.NIM))
	request.File = linkFile.SignedURL + fmt.Sprintf("?download=%s-%s-%d.pdf", sectionClass.Title, student.Name, student.NIM)

	customError = s.classRepo.InsertSubmissionStudent(c, request)
	if customError.Cause != nil {
		return customError
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) FetchSubmissionBySection(c context.Context, sectionClassId int, teacherId uuid.UUID) ([]*models.StudentSubmission, pkg.CustomError) {
	//isAuthorized, customError := s.classRepo.CheckStudentClassExists(c, teacherId, )

	submissions, customError := s.classRepo.GetSubmissionByClassSection(c, sectionClassId)
	if customError.Cause != nil {
		return nil, customError
	}

	return submissions, pkg.CustomError{}
}

func NewClassUsecase(classRepo repository.ClassRepository, studentRepo repository.StudentRepository) ClassUsecase {
	return &classUsecaseImpl{
		classRepo:   classRepo,
		studentRepo: studentRepo,
	}
}
