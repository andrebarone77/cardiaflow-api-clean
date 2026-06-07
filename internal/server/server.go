package server

import (
	"database/sql"

	"github.com/andrebarone77/cardiaflow-api/configs"
	"github.com/andrebarone77/cardiaflow-api/internal/handler"
	"github.com/andrebarone77/cardiaflow-api/internal/repository"
	"github.com/andrebarone77/cardiaflow-api/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	port string
	db   *sql.DB
}

func NewServer(dbase *sql.DB, cfg *configs.Config) *Server {
	return &Server{
		port: cfg.AppPort,
		db:   dbase,
	}
}

func (s *Server) Run() {
	r := gin.Default()

	userRepo := repository.NewUserRepository(s.db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	healthRecordTypeRepo := repository.NewHealthRecordTypeRepository(s.db)
	healthRecordTypeService := service.NewHealthRecordTypeService(healthRecordTypeRepo)
	healthRecordTypeHandler := handler.NewHealthRecordTypeHandler(healthRecordTypeService)

	healthRecordRepo := repository.NewHealthRecordRepository(s.db)
	healthRecordService := service.NewHealthRecordService(healthRecordRepo)
	healthRecordHandler := handler.NewHealthRecordHandler(healthRecordService)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	api := r.Group("/api")
	{
		api.POST("/users", userHandler.Create)
		api.GET("/users", userHandler.Get)
		api.GET("/users/:id", userHandler.GetById)
		api.DELETE("/users/:id", userHandler.Delete)
		api.PATCH("/users/:id", userHandler.Update)

		api.POST("/healthrecordtypes", healthRecordTypeHandler.Create)
		api.GET("/healthrecordtypes", healthRecordTypeHandler.GetAll)
		api.GET("/healthrecordtypes/:id", healthRecordTypeHandler.GetByID)
		api.GET("/healthrecordtypes/code/:code", healthRecordTypeHandler.GetByCode)
		api.DELETE("/healthrecordtypes", healthRecordTypeHandler.Delete)
		api.PATCH("/healthrecordtypes/:id", healthRecordTypeHandler.Update)

		api.POST("/healthrecord", healthRecordHandler.Create)
		api.GET("/healthrecord/:id", healthRecordHandler.GetByID)
		api.GET("/healthrecord/list", healthRecordHandler.ListByUserID)
		api.DELETE("/healthrecord", healthRecordHandler.Delete)
		api.PATCH("/healthrecord/:id", healthRecordHandler.Update)

		api.POST("/auth/login", authHandler.Login)

	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":" + s.port)
}
