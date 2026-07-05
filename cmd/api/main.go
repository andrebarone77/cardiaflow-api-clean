// @title           Cardiaflow API
// @version         1.0
// @description     REST API for managing users and health records.
// @termsOfService  https://github.com/andrebarone77/cardiaflow-api

// @contact.name   André Barone
// @contact.url    https://github.com/andrebarone77
// @contact.email  andrebarone77@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"github.com/andrebarone77/cardiaflow-api/configs"
	"github.com/andrebarone77/cardiaflow-api/internal/database"
	"github.com/andrebarone77/cardiaflow-api/internal/server"
)

func main() {
	cfg := configs.Load()
	db := database.New(cfg)
	srv := server.NewServer(db, cfg)

	defer db.Close()
	srv.Run()
}
