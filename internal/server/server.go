package server

import (
	"database/sql"

	"github.com/andrebarone77/cardiaflow-api/configs"
	"github.com/andrebarone77/cardiaflow-api/internal/auth"
	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/andrebarone77/cardiaflow-api/internal/handler"
	"github.com/andrebarone77/cardiaflow-api/internal/repository"
	"github.com/andrebarone77/cardiaflow-api/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/andrebarone77/cardiaflow-api/docs"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	login := r.Group("/api/auth")
	login.POST("/login", authHandler.Login)

	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	{

		api.GET("/healthrecord/:id", healthRecordHandler.GetByID)
		api.GET("/healthrecord/list", healthRecordHandler.ListByUserID)
		api.DELETE("/healthrecord", healthRecordHandler.Delete)
		api.PATCH("/healthrecord/:id", healthRecordHandler.Update)

	}

	healtRecordTypesHND := r.Group("/api")
	{
		healtRecordTypesHND.POST("/healthrecordtypes", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), healthRecordTypeHandler.Create)
		healtRecordTypesHND.GET("/healthrecordtypes", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser), healthRecordTypeHandler.GetAll)
		healtRecordTypesHND.GET("/healthrecordtypes/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser), healthRecordTypeHandler.GetByID)
		healtRecordTypesHND.GET("/healthrecordtypes/code/:code", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser), healthRecordTypeHandler.GetByCode)
		healtRecordTypesHND.DELETE("/healthrecordtypes", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), healthRecordTypeHandler.Delete)
		healtRecordTypesHND.PATCH("/healthrecordtypes/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), healthRecordTypeHandler.Update)
	}

	usersHND := r.Group("/api")
	{
		usersHND.GET("/users", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser), userHandler.Get)
		usersHND.GET("/users/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser), userHandler.GetById)
		usersHND.DELETE("/users/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), userHandler.Delete)
		usersHND.PATCH("/users/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), userHandler.Update)
		usersHND.POST("/users", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), userHandler.Create)
	}

	create := r.Group("/api")
	{

		create.POST("/healthrecord", healthRecordHandler.Create)
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":" + s.port)
}
