package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rifkhia/lms-remake/internal"
	"github.com/rifkhia/lms-remake/internal/delivery/handler"
	classRepository "github.com/rifkhia/lms-remake/internal/repository/class_impl"
	studentRepository "github.com/rifkhia/lms-remake/internal/repository/student_impl"
	teacherRepository "github.com/rifkhia/lms-remake/internal/repository/teacher_impl"
	"github.com/rifkhia/lms-remake/internal/usecase"
	"github.com/spf13/viper"
)

func initViperConfig() {
	viper.SetConfigName("app")

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func main() {
	initViperConfig()

	database := internal.ConnectDatabase()

	studentRepository := studentRepository.NewStudentRepository(database)
	classRepository := classRepository.NewClassRepository(database)
	teacherRepository := teacherRepository.NewTeacherRepository(database)
	studentUsecase := usecase.NewStudentUsecase(studentRepository)
	classUsecase := usecase.NewClassUsecase(classRepository)
	teacherUsecase := usecase.NewTeacherUsecase(teacherRepository)
	studentHandler := handler.NewStudentHandler(studentUsecase)
	classHandler := handler.NewClassHandler(classUsecase)
	teacherHandler := handler.NewTeacherHandler(teacherUsecase)
	app := fiber.New()

	app.Use(logger.New())

	studentHandler.Route(app)
	classHandler.Route(app)
	teacherHandler.Route(app)

	app.Listen(":8081")
}
