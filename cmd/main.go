package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rifkhia/lms-remake/internal"
	"github.com/rifkhia/lms-remake/internal/delivery/handler"
	"github.com/rifkhia/lms-remake/internal/repository"
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

//func initSupabaseStorage() {
//	storageClient := storage_go.NewClient("https://sfiaqorbwfekbitsqvsf.supabase.co/storage/v1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InNmaWFxb3Jid2Zla2JpdHNxdnNmIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MDcxNTYzNzMsImV4cCI6MjAyMjczMjM3M30.-_P0sYfzdbgeoHvXAm0AmFpUiv3n7feqQWQq1N22MoA", nil)
//}

func main() {
	initViperConfig()

	database := internal.ConnectDatabase()

	studentRepository := repository.NewStudentRepository(database)
	classRepository := repository.NewClassRepository(database)
	teacherRepository := repository.NewTeacherRepository(database)
	studentUsecase := usecase.NewStudentUsecase(studentRepository)
	classUsecase := usecase.NewClassUsecase(classRepository, studentRepository)
	teacherUsecase := usecase.NewTeacherUsecase(teacherRepository)
	studentHandler := handler.NewStudentHandler(studentUsecase)
	classHandler := handler.NewClassHandler(classUsecase, studentUsecase)
	teacherHandler := handler.NewTeacherHandler(teacherUsecase)
	app := fiber.New()

	app.Use(logger.New())

	studentHandler.Route(app)
	classHandler.Route(app)
	teacherHandler.Route(app)

	app.Listen(":8081")
}
