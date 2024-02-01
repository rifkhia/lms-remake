package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/delivery/middleware"
	"github.com/rifkhia/lms-remake/internal/usecase"
	"strconv"
)

type ClassHandlerImpl struct {
	classUsecase usecase.ClassUsecase
}

func (handler ClassHandlerImpl) Route(app *fiber.App) {
	app.Get("/v1/class/:id", middleware.JWTGuardAll, handler.FetchClassById)
	app.Get("/v1/class", middleware.JWTGuardAll, handler.FetchClassByName)
	app.Post("/v1/class/join", middleware.JWTGuardStudent, handler.StudentJoinClass)
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

	classResult, err := handler.classUsecase.FetchClassById(c.Context(), classId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
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

	classResult, err := handler.classUsecase.FetchClassByName(c.Context(), param)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Success get class by name: %s", param),
		"data":    classResult,
	})
}

func (handler *ClassHandlerImpl) StudentJoinClass(c *fiber.Ctx) error {
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
		Id  int    `json:"id"`
		Key string `json:"key"`
	}{}
	c.BodyParser(&request)
	err = handler.classUsecase.JoinClass(c.Context(), parsedStudentId, request.Id, request.Key)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success join class",
	})
}

func NewClassHandler(classUsecase usecase.ClassUsecase) *ClassHandlerImpl {
	return &ClassHandlerImpl{
		classUsecase: classUsecase,
	}
}
