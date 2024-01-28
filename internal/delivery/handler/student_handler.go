package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/delivery/middleware"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/usecase"
)

func (handler StudentHandlerImpl) Route(app *fiber.App) {
	app.Get("/v1/student", middleware.JWTGuardStudent, handler.FetchStudent)
	app.Post("/v1/student", handler.RegisterStudent)
	app.Post("/v1/student/login", handler.LoginStudent)
}

type StudentHandlerImpl struct {
	studentUsecase usecase.StudentUsecase
}

func (handler *StudentHandlerImpl) FetchStudent(c *fiber.Ctx) error {
	var studentResult interface{}
	var err error

	if c.Query("id") != "" {
		param := c.Query("id")
		id, err := uuid.Parse(param)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Wrong UUID"))
		}

		studentResult, err = handler.studentUsecase.FetchStudentById(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
				"message": err,
			})
		}
	} else if c.Query("name") != "" {
		param := c.Query("name")

		studentResult, err = handler.studentUsecase.FetchStudentByName(c.Context(), param)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
				"message": err,
			})
		}

	} else {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"message": "Insert id or name!",
		})
	}

	return c.JSON(map[string]interface{}{
		"message": "Success fetching student data",
		"data":    studentResult,
	})
}

func (handler *StudentHandlerImpl) RegisterStudent(c *fiber.Ctx) error {
	var request models.LoginRegisterRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"message": err,
		})
	}

	data, err := handler.studentUsecase.Register(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"message": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]interface{}{
		"message": "Success create student",
		"data":    data,
	})
}

func (handler *StudentHandlerImpl) LoginStudent(c *fiber.Ctx) error {
	var request models.LoginRegisterRequest
	var data interface{}
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"message": err.Error(),
		})
	}

	data, err = handler.studentUsecase.Login(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"messsage": "Login success",
		"data":     data,
	})
}
func NewStudentHandler(studentUsecase usecase.StudentUsecase) *StudentHandlerImpl {
	return &StudentHandlerImpl{
		studentUsecase: studentUsecase,
	}
}
