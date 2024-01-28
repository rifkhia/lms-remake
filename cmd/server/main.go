package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rifkhia/lms-remake/internal"
	"github.com/rifkhia/lms-remake/internal/delivery/handler"
	studentRepository "github.com/rifkhia/lms-remake/internal/repository/student_impl"
	"github.com/rifkhia/lms-remake/internal/usecase"
)

func main() {
	godotenv.Load(".env")

	database := internal.ConnectDatabase()

	studentRepository := studentRepository.NewStudentRepository(database)

	studentUsecase := usecase.NewStudentUsecase(studentRepository)

	studentHandler := handler.NewStudentHandler(studentUsecase)

	app := fiber.New()

	studentHandler.Route(app)

	app.Listen(":8080")
}
