package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/delivery/middleware"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/usecase"
	"github.com/rifkhia/lms-remake/internal/utils"
)

func (handler StudentHandlerImpl) Route(app *fiber.App) {
	app.Get("/v1/student/profile", middleware.JWTGuardStudent, handler.FetchStudentById)
	app.Post("/v1/student/register", handler.RegisterStudent)
	app.Post("/v1/student/login", handler.LoginStudent)
	app.Delete("/v1/student", middleware.JWTGuardStudent, handler.DeleteStudent)
}

type StudentHandlerImpl struct {
	studentUsecase usecase.StudentUsecase
}

func (handler *StudentHandlerImpl) FetchStudentById(c *fiber.Ctx) error {
	id, err := middleware.GetIdFromToken(c)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	parseId, err := uuid.Parse(id)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	studentResult, customError := handler.studentUsecase.FetchStudentById(c.Context(), parseId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.JSON(map[string]interface{}{
		"message": "Success fetching student data",
		"data":    studentResult,
	})
}

func (handler *StudentHandlerImpl) FetchStudentByName(c *fiber.Ctx) error {
	param := c.Query("name")
	if param == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "name cannot be blank!",
		})
	}

	studentResult, customError := handler.studentUsecase.FetchStudentByName(c.Context(), param)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.JSON(map[string]interface{}{
		"message": "Success fetching student data",
		"data":    studentResult,
	})
}

func (handler *StudentHandlerImpl) RegisterStudent(c *fiber.Ctx) error {
	var request models.StudentRegisterRequest
	err := c.BodyParser(&request)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	data, customError := handler.studentUsecase.Register(c.Context(), &request)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"message": "Success create student",
		"data":    data,
	})
}

func (handler *StudentHandlerImpl) LoginStudent(c *fiber.Ctx) error {
	var request models.StudentLoginRequest
	var data interface{}
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"message": err.Error(),
		})
	}

	data, customError := handler.studentUsecase.Login(c.Context(), &request)
	if customError.Error() != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"messsage": "Login success",
		"data":     data,
	})
}

func (handler *StudentHandlerImpl) DeleteStudent(c *fiber.Ctx) error {
	id, err := middleware.GetIdFromToken(c)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	parseId, err := uuid.Parse(id)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   err,
			Service: utils.HANDLER_SERVICE,
		}
		return c.Status(customError.Code).JSON(customError.Error())
	}

	customError := handler.studentUsecase.DeleteStudent(c.Context(), parseId)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account deleted",
	})
}

func NewStudentHandler(studentUsecase usecase.StudentUsecase) *StudentHandlerImpl {
	return &StudentHandlerImpl{
		studentUsecase: studentUsecase,
	}
}
