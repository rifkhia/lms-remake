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
	"github.com/supabase-community/storage-go"
	"strconv"
)

type ClassHandlerImpl struct {
	classUsecase usecase.ClassUsecase
}

func (handler ClassHandlerImpl) Route(app *fiber.App) {
	app.Get("/v1/class/:id", middleware.JWTGuardAll, handler.FetchClassById)
	app.Get("/v1/class", middleware.JWTGuardAll, handler.FetchClassByName)
	app.Post("v1/class", middleware.JWTGuardTeacher, handler.CreateClass)
	app.Post("/v1/class/:id/join", middleware.JWTGuardStudent, handler.StudentJoinClass)
	app.Post("v1/class/:id/section", middleware.JWTGuardTeacher, handler.CreateClassSection)
	app.Post("v1/class/:id/section/:section_id/submissions", middleware.JWTGuardTeacher, handler.AddSubmissionsTeacher)
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

	teacherId, err := middleware.GetIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	parsedTeacherId, err := uuid.Parse(teacherId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	classId := c.Params("id")
	if classId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "id cannot be blank!",
		})
	}

	intClassId, err := strconv.Atoi(classId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"messsage": err.Error(),
		})
	}

	class, customError := handler.classUsecase.FetchClassById(c.Context(), intClassId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	if class.TeacherId != parsedTeacherId {
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

	var exists = false
	for _, v := range class.ClassSection {
		if intSectionClassId == v.ID {
			exists = true
		}
	}

	if !exists {
		customError := pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("no section found in this class"),
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

	storageGo := storage_go.NewClient("https://sfiaqorbwfekbitsqvsf.supabase.co/storage/v1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InNmaWFxb3Jid2Zla2JpdHNxdnNmIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTcwNzE1NjM3MywiZXhwIjoyMDIyNzMyMzczfQ.eLQlf2WyPRZLdgvtrIfN4bc_veo2kZiN9M1DS6iX2h0", nil)
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
	linkFile.SignedURL = linkFile.SignedURL + fmt.Sprintf("?download=%s.pdf", request.Title)
	return c.Status(fiber.StatusOK).JSON(linkFile)

}

func NewClassHandler(classUsecase usecase.ClassUsecase) *ClassHandlerImpl {
	return &ClassHandlerImpl{
		classUsecase: classUsecase,
	}
}
