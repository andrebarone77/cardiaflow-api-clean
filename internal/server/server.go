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

		api.GET("/healthrecordtypes", healthRecordTypeHandler.GetAll)
		api.GET("/healthrecordtypes/:id", healthRecordTypeHandler.GetByID)
		api.GET("/healthrecordtypes/code/:code", healthRecordTypeHandler.GetByCode)
		api.DELETE("/healthrecordtypes", healthRecordTypeHandler.Delete)
		api.PATCH("/healthrecordtypes/:id", healthRecordTypeHandler.Update)

		api.GET("/healthrecord/:id", healthRecordHandler.GetByID)
		api.GET("/healthrecord/list", healthRecordHandler.ListByUserID)
		api.DELETE("/healthrecord", healthRecordHandler.Delete)
		api.PATCH("/healthrecord/:id", healthRecordHandler.Update)

	}

	users := r.Group("/api")
	{
		users.GET("/users", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser), userHandler.Get)
		users.GET("/users/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser), userHandler.GetById)
		users.DELETE("/users/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), userHandler.Delete)
		users.PATCH("/users/:id", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), userHandler.Update)
		users.POST("/users", auth.AuthMiddleware(), auth.RequireRoles(domain.RoleAdmin, domain.RoleManager), userHandler.Create)
	}

	create := r.Group("/api")
	{

		create.POST("/healthrecordtypes", healthRecordTypeHandler.Create)
		create.POST("/healthrecord", healthRecordHandler.Create)
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":" + s.port)
}
