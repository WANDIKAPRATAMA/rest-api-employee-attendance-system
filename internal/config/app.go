package config

import (
	controller "employee-attendance-system/internal/controllers"
	"employee-attendance-system/internal/middleware"
	"employee-attendance-system/internal/repository"
	route "employee-attendance-system/internal/route"
	"employee-attendance-system/internal/usecase"
	utils "employee-attendance-system/internal/util"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Viper    *viper.Viper
}

func NewAppConfig(config *AppConfig) {

	jwtUtils := utils.NewJWTCfg(config.Viper)

	userRepo := repository.NewUserRepository(config.DB, config.Log)
	authUseCase := usecase.NewAuthUseCase(userRepo, config.Log, config.Validate, config.Viper, jwtUtils)
	authController := controller.NewAuthController(authUseCase, config.Log, config.Validate)
	authMiddleware := middleware.NewAuth(authUseCase, config.Log, config.Viper, jwtUtils)

	deptRepo := repository.NewDepartmentRepository(config.DB, config.Log)
	deptUseCase := usecase.NewDepartmentUseCase(deptRepo, config.Log, config.Validate)
	deptController := controller.NewDepartmentController(deptUseCase, config.Log, config.Validate)

	attRepo := repository.NewAttendanceRepository(config.DB, config.Log)
	attUseCase := usecase.NewAttendanceUseCase(attRepo, userRepo, config.Log, config.Validate) // Reuse profileRepo
	attController := controller.NewAttendanceController(attUseCase, config.Log, config.Validate)

	authRoutesConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}

	profileRoutesConfig := route.ProfileRouteConfig{
		App:            config.App,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}

	attRoutesConfig := route.AttendanceRouteConfig{
		App:                  config.App,
		AttendanceController: attController,
		AuthMiddleware:       authMiddleware,
	}
	deptRoutesConfig := route.DepartmentRouteConfig{
		App:                  config.App,
		DepartmentController: deptController,
		AuthMiddleware:       authMiddleware,
	}
	authRoutesConfig.Setup()
	profileRoutesConfig.Setup()
	deptRoutesConfig.Setup()
	attRoutesConfig.Setup()
	config.Log.Info("Server starting on :8080")
	if err := config.App.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
