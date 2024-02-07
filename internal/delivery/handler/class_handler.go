package handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/delivery/middleware"
	"github.com/rifkhia/lms-remake/internal/dto"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/usecase"
	"github.com/rifkhia/lms-remake/internal/utils"
	"github.com/spf13/viper"
	"github.com/supabase-community/storage-go"
	"strconv"
)

type ClassHandlerImpl struct {
	classUsecase   usecase.ClassUsecase
	studentUsecase usecase.StudentUsecase
}

func (handler ClassHandlerImpl) Route(app *fiber.App) {
	app.Get("/v1/class/:id", middleware.JWTGuardAll, handler.FetchClassById)
	app.Get("/v1/class", middleware.JWTGuardAll, handler.FetchClassByName)
	app.Post("v1/class", middleware.JWTGuardTeacher, handler.CreateClass)
	app.Post("/v1/class/:id/join", middleware.JWTGuardStudent, handler.StudentJoinClass)
	app.Post("v1/class/:id/section", middleware.JWTGuardTeacher, handler.CreateClassSection)
	app.Post("v1/class/section/:section_id/submissions", middleware.JWTGuardTeacher, handler.AddSubmissionsTeacher)
	app.Get("v1/class/section/:section_id/submissions", middleware.JWTGuardTeacher, handler.FetchSubmission)
	app.Post("v1/class/section/:section_id", middleware.JWTGuardStudent, handler.AddSubmissionsStudent)
}

func (handler *ClassHandlerImpl) FetchClassById(c *fiber.Ctx) error {
	param := c.Params("id")
	if param == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id cannot be blank!",
		})
	}

	classId, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "input integer only for class id",
		})
	}

	classResult, customError := handler.classUsecase.FetchClassById(c.Context(), classId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Success get class by id: %d", classId),
		"data":    classResult,
	})
}

func (handler *ClassHandlerImpl) FetchClassByName(c *fiber.Ctx) error {
	param := c.Query("name")
	if param == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "class name cannot be blank!",
		})
	}

	classResult, customError := handler.classUsecase.FetchClassByName(c.Context(), param)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Success get class by name: %s", param),
		"data":    classResult,
	})
}

func (handler *ClassHandlerImpl) CreateClass(c *fiber.Ctx) error {
	var request *dto.ClassCreate
	teacherId, err := middleware.GetIdFromToken(c)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	parsedTeacherId, err := uuid.Parse(teacherId)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	err = c.BodyParser(&request)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	customError := handler.classUsecase.CreateClass(c.Context(), request, parsedTeacherId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "class created",
	})
}

func (handler *ClassHandlerImpl) StudentJoinClass(c *fiber.Ctx) error {
	param := c.Params("id")
	if param == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "class id cannot be blank!",
		})
	}

	intId, _ := strconv.Atoi(param)

	studentId, err := middleware.GetIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	parsedStudentId, err := uuid.Parse(studentId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}
	request := struct {
		Key string `json:"key"`
	}{}
	err = c.BodyParser(&request)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.UNPROCESSABLE_ENTITY,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	customError := handler.classUsecase.JoinClass(c.Context(), parsedStudentId, intId, request.Key)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success join class",
	})
}

func (handler *ClassHandlerImpl) CreateClassSection(c *fiber.Ctx) error {
	var request models.SectionClass
	teacherId, err := middleware.GetIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	param := c.Params("id")
	if param == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id cannot be blank!",
		})
	}

	err = c.BodyParser(&request)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	request.ClassId, _ = strconv.Atoi(param)

	class, customError := handler.classUsecase.FetchClassById(c.Context(), request.ClassId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	parseTeacherId, err := uuid.Parse(teacherId)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	if class.TeacherId != parseTeacherId {
		customError = pkg.CustomError{
			Code:    utils.FORBIDDEN,
			Cause:   errors.New("you don't have authority to this class"),
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	customError = handler.classUsecase.CreateSectionClass(c.Context(), &request)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("section class %s created", request.Title),
	})
}

func (handler *ClassHandlerImpl) AddSubmissionsTeacher(c *fiber.Ctx) error {
	var request models.Submission

	Id, err := middleware.GetIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	parsedId, err := uuid.Parse(Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	sectionClassId := c.Params("section_id")
	if sectionClassId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id cannot be blank!",
		})
	}

	intSectionClassId, err := strconv.Atoi(sectionClassId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	sectionClass, customError := handler.classUsecase.FetchSectionClassById(c.Context(), intSectionClassId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	class, customError := handler.classUsecase.FetchClassById(c.Context(), sectionClass.ClassId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	if class.TeacherId != parsedId {
		customError = pkg.CustomError{
			Code:    utils.FORBIDDEN,
			Cause:   errors.New("you don't have authority to this class"),
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	err = c.BodyParser(&request)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	request.ClassSectionId = intSectionClassId

	file, err := c.FormFile("file")
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	parsedFile, err := file.Open()
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	storageGo := storage_go.NewClient(viper.GetString("SUPABASE_URL"), viper.GetString("SUPABASE_TOKEN"), nil)
	_, err = storageGo.UploadFile("submissions_teacher", fmt.Sprintf("class_section/%d/%s.pdf", request.ClassSectionId, request.Title), parsedFile)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	linkFile := storageGo.GetPublicUrl("submissions_teacher", fmt.Sprintf("class_section/%d/%s.pdf", request.ClassSectionId, request.Title))
	request.File = linkFile.SignedURL + fmt.Sprintf("?download=%s.pdf", request.Title)

	customError = handler.classUsecase.AddSubmissionTeacher(c.Context(), &request)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file uploaded",
		"link":    request.File,
	})
}

func (handler *ClassHandlerImpl) AddSubmissionsStudent(c *fiber.Ctx) error {
	var request dto.StudentSubmissionRequest

	Id, err := middleware.GetIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	parsedId, err := uuid.Parse(Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	sectionClassId := c.Params("section_id")
	if sectionClassId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id cannot be blank!",
		})
	}

	intSectionClassId, err := strconv.Atoi(sectionClassId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	err = c.BodyParser(&request)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	request = dto.StudentSubmissionRequest{
		ID:             parsedId,
		ClassSectionId: intSectionClassId,
	}

	customError := handler.classUsecase.AddSubmissionStudent(c.Context(), &request, file)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file uploaded",
		"link":    request.File,
	})
}

func (handler *ClassHandlerImpl) FetchSubmission(c *fiber.Ctx) error {
	Id, err := middleware.GetIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	parsedId, err := uuid.Parse(Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	sectionClassId := c.Params("section_id")
	if sectionClassId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id cannot be blank!",
		})
	}

	intSectionClassId, err := strconv.Atoi(sectionClassId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	submissions, customError := handler.classUsecase.FetchSubmissionBySection(c.Context(), intSectionClassId, parsedId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success getting submissions",
		"data":    submissions,
	})
}

func NewClassHandler(classUsecase usecase.ClassUsecase, studentUsecase usecase.StudentUsecase) *ClassHandlerImpl {
	return &ClassHandlerImpl{
		classUsecase:   classUsecase,
		studentUsecase: studentUsecase,
	}
}
