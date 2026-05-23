#!/bin/bash
mkdir -p cmd/api
touch cmd/api/main.go

mkdir -p internal/app
touch internal/app/app.go
touch internal/app/router.go

mkdir -p internal/domain
touch internal/domain/user.go
touch internal/domain/health_record.go
touch internal/domain/goal.go

mkdir -p internal/service
touch internal/service/user_service.go
touch internal/service/health_service.go
touch internal/service/goal_service.go

mkdir -p internal/repository
touch internal/repository/user_repository.go
touch internal/repository/health_repository.go
touch internal/repository/goal_repository.go

mkdir -p internal/handler
touch internal/handler/user_handler.go
touch internal/handler/health_handler.go
touch internal/handler/goal_handler.go

mkdir -p internal/midleware
touch internal/midleware/auth.go

mkdir -p pkg/utils
touch pkg/utils/hash.go
touch pkg/utils/jwt.go

mkdir configs
touch configs/config.go

touch go.mod
touch README.md