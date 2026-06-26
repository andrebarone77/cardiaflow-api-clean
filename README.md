# Cardiaflow API

Cardiaflow is a REST API for storing and managing personal health records such as weight, blood pressure, glucose levels and other health metrics.


## Development:
- Golang with [Gin Web Framework](https://gin-gonic.com/)
- PostgreSQL
- JWT authentication


## Configuração
### Install go language  
    Linux Ubuntu:   
    sudo snap install go --classic  

    Windows:  
    Using Power Shell or Command prompt execute
    winget install GoLang.Go  
    agree with the terms (if you feel comfortable :D )  
    give permission for the installer to go on  
    restart the terminal (Power shell or Command Prompt) 
    Test with go version.


### Database Setup

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);  
```
```sql
CREATE TABLE health_record_types (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name text NOT NULL CHECK (name = LOWER(TRIM(name))),
	code text NOT NULL CHECK (code = LOWER(TRIM (code))),
	unit text CHECK (unit= LOWER(TRIM (unit))),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
CREATE UNIQUE INDEX health_record_types_code_key
ON health_record_types (LOWER(TRIM(code)));

CREATE UNIQUE INDEX users_email_unique
ON users (LOWER(email));

CREATE EXTENSION IF NOT EXISTS pgcrypto;
```
```sql
CREATE TABLE health_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    health_record_type_id UUID NOT NULL,
    value NUMERIC(12,4) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_health_records_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
   CONSTRAINT fk_health_records_type
        FOREIGN KEY (health_record_type_id)
        REFERENCES health_record_types(id)
        ON DELETE RESTRICT,
    CONSTRAINT health_records_value_check
        CHECK (value >= 0),
    CONSTRAINT health_records_notes_check
        CHECK (notes IS NULL OR length(trim(notes)) > 0)
);    
``` 
### Enviroment Example (.env)
```
APP_PORT=8080

DB_HOST=<host_ip or name>
DB_PORT=5433
DB_USER=<dbuser>
DB_PASSWORD=<dbpassword>
DB_NAME=cardiaflow
DB_SSLMODE=disable

JWT_SECRET=<jwtsecret>
JWT_EXPIRES_IN=2h
```

## Main Endpoints

    [POST] /api/auth/login - login

    [POST] /users - create user
    [GET] /users - get user by email [jwt]
    [GET] /users/:id - get user by id [jwt]
    [DELETE] /users/:id - delete user [jwt]
    [PATCH] /users/:id - update user [jwt]

    [POST] /healthrecordtypes - create health record type
    [GET] /healthrecordtypes - get all health record types [jwt]
    [GET] /healthrecordtypes/:id - get health record type by id [jwt]
    [GET] /healthrecordtypes/code/:code - get health record type by code [jwt]
    [DELETE] /healthrecordtypes- delete health record type [jwt]
    [PATCH] /healthrecordtypes/:id - update health record type [jwt]

    [POST] /healthrecord - create health record 
    [GET] /healthrecord/:id - get health record type by id [jwt]
    [GET] /healthrecord/list- get health record list by user id [jwt]
    [DELETE] /healthrecord- delete health record [jwt]
    [PATCH] /healthrecord/:id - update health record [jwt]


## Roadmap
    - [ ] User Role and permissions
    - [ ] Unit Tests
    - [ ] Swagger documentation
    - [ ] Frontend application

## Useful
### To Run Tests and show coverage by function:
```
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
```
### To Run Tests and show covarage in HTML:
```
go test ./internal/handler/ -coverprofile=coverage.out && go tool cover -html=coverage.out
```
