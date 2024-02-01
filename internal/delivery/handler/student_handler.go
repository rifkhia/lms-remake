package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/delivery/middleware"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/usecase"
)

func (handler StudentHandlerImpl) Route(app *fiber.App) {
	app.Get("/v1/student", middleware.JWTGuardStudent, handler.FetchStudentById)
	app.Post("/v1/student", handler.RegisterStudent)
	app.Post("/v1/student/login", handler.LoginStudent)
}

type StudentHandlerImpl struct {
	studentUsecase usecase.StudentUsecase
}

func (handler *StudentHandlerImpl) FetchStudentById(c *fiber.Ctx) error {
	id, err := middleware.GetIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err,
		})
	}

	parseId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messsage": err,
		})
	}

	studentResult, err := handler.studentUsecase.FetchStudentById(c.Context(), parseId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"message": err,
		})
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

	studentResult, err := handler.studentUsecase.FetchStudentByName(c.Context(), param)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
			"message": err,
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
		if err.Error() == "Email not found!" {
			return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
				"message": err.Error(),
			})
		}
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
