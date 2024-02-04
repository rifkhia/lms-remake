package handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/delivery/middleware"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/usecase"
	"github.com/rifkhia/lms-remake/internal/utils"
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
	var request *models.ClassCreate
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
		return c.Status(customError.Code).JSON(fiber.Map{
			"1": class,
			"2": parseTeacherId,
			"3": request.ClassId,
		})
	}

	customError = handler.classUsecase.CreateSectionClass(c.Context(), &request)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("section class %s created", request.Title),
	})
}

func NewClassHandler(classUsecase usecase.ClassUsecase) *ClassHandlerImpl {
	return &ClassHandlerImpl{
		classUsecase: classUsecase,
	}
}
