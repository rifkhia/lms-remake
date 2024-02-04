package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/usecase"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type TeacherHandlerImpl struct {
	teacherUsecase usecase.TeacherUsecase
}

func (handler TeacherHandlerImpl) Route(app *fiber.App) {
	app.Post("/v1/teacher/login", handler.LoginTeacher)
	app.Post("/v1/teacher/register", handler.RegisterTeacher)

}

func (handler *TeacherHandlerImpl) LoginTeacher(c *fiber.Ctx) error {
	var request models.TeacherLoginRequest
	var data interface{}

	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if request.Email == "" && request.NPM == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email or npm must be filled",
		})
	}

	if request.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "password can't be blank",
		})
	}

	data, customError := handler.teacherUsecase.LoginTeacher(c.Context(), &request)
	if customError.Cause != nil {
		return c.Status(customError.Code).JSON(customError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login success",
		"data":    data,
	})
}

func (handler *TeacherHandlerImpl) RegisterTeacher(c *fiber.Ctx) error {
	var request models.TeacherRegisterRequest
	err := c.BodyParser(&request)
	if err != nil {
		customError := pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: "Handler",
		}
		return c.Status(utils.INTERNAL_SERVER_ERROR).JSON(customError.Error())
	}

	data, custErr := handler.teacherUsecase.RegisterTeacher(c.Context(), &request)
	if custErr.Cause != nil {
		return c.Status(custErr.Code).JSON(custErr.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "register success",
		"data":    data,
	})
}

func NewTeacherHandler(teacherUsecase usecase.TeacherUsecase) *TeacherHandlerImpl {
	return &TeacherHandlerImpl{
		teacherUsecase: teacherUsecase,
	}
}
