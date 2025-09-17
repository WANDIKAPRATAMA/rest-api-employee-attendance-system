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
	config.App.Use(middleware.SetupCORS())
	config.App.Use(middleware.SetupRateLimiter())

	jwtUtils := utils.NewJWTCfg(config.Viper)

	userRepo := repository.NewUserRepository(config.DB, config.Log)
	authUseCase := usecase.NewAuthUseCase(userRepo, config.Log, config.Validate, config.Viper, jwtUtils)
	authController := controller.NewAuthController(authUseCase, config.Log, config.Validate)
	authMiddleware := middleware.NewAuth(authUseCase, config.Log, config.Viper, jwtUtils)

	userUseCase := usecase.NewUserUseCase(userRepo, config.Log, config.Validate)
	userController := controller.NewUserController(userUseCase, config.Log, config.Validate)

	deptRepo := repository.NewDepartmentRepository(config.DB, config.Log)
	deptUseCase := usecase.NewDepartmentUseCase(deptRepo, config.Log, config.Validate, userRepo)
	deptController := controller.NewDepartmentController(deptUseCase, config.Log, config.Validate)

	attRepo := repository.NewAttendanceRepository(config.DB, config.Log)
	attUseCase := usecase.NewAttendanceUseCase(attRepo, userRepo, deptRepo, config.Log, config.Validate) // Reuse profileRepo
	attController := controller.NewAttendanceController(attUseCase, config.Log, config.Validate)

	authRoutesConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}

	profileRoutesConfig := route.UserRouteConfig{
		App:            config.App,
		UserController: userController,
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
