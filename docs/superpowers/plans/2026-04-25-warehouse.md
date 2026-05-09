# Warehouse Management System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a warehouse management system with inventory, orders, audit logging, and RBAC permissions.

**Architecture:** Go backend with Gin + Bun ORM, React frontend with TypeScript + Ant Design. Embedded SPA deployment with static files compiled into the Go binary.

**Tech Stack:** Go 1.21+, Gin, Bun ORM, MySQL/MariaDB, React 18, TypeScript, Vite, Ant Design 5.x

---

## Phase 1: Project Setup & Infrastructure

### Task 1.1: Initialize Go Project

- [ ] **Step 1: Initialize Go module**

```bash
cd /home/zzf/projects/goinvent
mkdir -p warehouse/cmd/server
mkdir -p warehouse/internal/{config,model,repository,service,handler,middleware,pkg/{response,errors,jwt,password},router}
mkdir -p warehouse/pkg/logger
mkdir -p warehouse/migrations
mkdir -p warehouse/config
cd warehouse && go mod init warehouse
```

- [ ] **Step 2: Create Makefile**

Create `warehouse/Makefile`:
```makefile
.PHONY: build run test clean dev

build:
	go build -o bin/warehouse ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin/

dev:
	go run ./cmd/server --config config/config.yaml

frontend-install:
	cd web && npm install

frontend-build:
	cd web && npm run build

all: build
```

- [ ] **Step 3: Commit**

```bash
git add . && git commit -m "feat: initialize go project"
```

---

### Task 1.2: Configuration Management

- [ ] **Step 1: Install dependencies**

```bash
cd /home/zzf/projects/goinvent/warehouse
go get gopkg.in/yaml.v3
go get github.com/gin-gonic/gin
go get github.com/uptrace/bun
go get github.com/uptrace/bun/driver/mysql
go get github.com/uptrace/bun/dialect/mysqldialect
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

- [ ] **Step 2: Create config struct**

Create `internal/config/config.go`:
```go
package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Name         string `yaml:"name"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
	File   string `yaml:"file"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

- [ ] **Step 3: Create config file**

Create `config/config.yaml`:
```yaml
server:
  port: 8080
  mode: debug

database:
  driver: mysql
  host: localhost
  port: 3306
  name: warehouse
  user: root
  password: ""
  max_open_conns: 100
  max_idle_conns: 10

jwt:
  secret: your-secret-key-change-in-production
  expire: 24h

log:
  level: debug
  output: stdout
  file: ""
```

- [ ] **Step 4: Commit**

```bash
git add . && git commit -m "feat: add configuration management"
```

---

### Task 1.3: Database Connection and Base Model

- [ ] **Step 1: Create database connection**

Create `internal/repository/db.go`:
```go
package repository

import (
	"database/sql"
	"fmt"

	"warehouse/internal/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/driver/mysqldriver"
)

func NewDB(cfg *config.DatabaseConfig) (*bun.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	sqldb, err := sql.Open(mysqldriver.DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)

	db := bun.NewDB(sqldb, mysqldialect.New())
	return db, nil
}
```

- [ ] **Step 2: Create base model**

Create `internal/model/base.go`:
```go
package model

import "time"

type BaseModel struct {
	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"created_at"`
	CreatedBy int64     `bun:"created_by,notnull" json:"created_by"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updated_at"`
	UpdatedBy int64     `bun:"updated_by,notnull" json:"updated_by"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

func (m *BaseModel) BeforeCreate(userID int64) {
	now := time.Now()
	m.CreatedAt = now
	m.CreatedBy = userID
	m.UpdatedAt = now
	m.UpdatedBy = userID
}

func (m *BaseModel) BeforeUpdate(userID int64) {
	m.UpdatedAt = time.Now()
	m.UpdatedBy = userID
}
```

- [ ] **Step 3: Create error definitions**

Create `internal/pkg/errors/errors.go`:
```go
package errors

import "errors"

type AppError struct {
	Code    int
	Message string
	Detail  string
}

func (e *AppError) Error() string {
	if e.Detail != "" {
		return e.Message + ": " + e.Detail
	}
	return e.Message
}

func NewAppError(code int, message, detail string) *AppError {
	return &AppError{Code: code, Message: message, Detail: detail}
}

const (
	CodeSuccess           = 0
	CodeBadRequest        = 400
	CodeUnauthorized      = 401
	CodeForbidden         = 403
	CodeNotFound          = 404
	CodeInternalError     = 500
	CodeUserNotFound      = 1001
	CodeInvalidPassword   = 1002
	CodeDuplicateEntry    = 1005
	CodeRecordNotFound    = 1006
	CodeInsufficientStock = 1004
)

var (
	ErrUserNotFound      = NewAppError(CodeUserNotFound, "用户不存在", "")
	ErrInvalidPassword   = NewAppError(CodeInvalidPassword, "密码错误", "")
	ErrDuplicateEntry    = NewAppError(CodeDuplicateEntry, "记录已存在", "")
	ErrRecordNotFound    = NewAppError(CodeRecordNotFound, "记录不存在", "")
	ErrInsufficientStock = NewAppError(CodeInsufficientStock, "库存不足", "")
)

func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return NewAppError(CodeInternalError, "内部错误", err.Error())
}
```

- [ ] **Step 4: Create response helpers**

Create `internal/pkg/response/response.go`:
```go
package response

import (
	"net/http"

	"warehouse/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithPage(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: "success",
		Data: PageData{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func Error(c *gin.Context, err error) {
	appErr := errors.GetAppError(err)
	status := http.StatusBadRequest
	switch appErr.Code {
	case errors.CodeUnauthorized:
		status = http.StatusUnauthorized
	case errors.CodeForbidden:
		status = http.StatusForbidden
	case errors.CodeNotFound, errors.CodeRecordNotFound:
		status = http.StatusNotFound
	case errors.CodeInternalError:
		status = http.StatusInternalServerError
	}
	c.JSON(status, Response{
		Code:    appErr.Code,
		Message: appErr.Message,
		Data:    appErr.Detail,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    errors.CodeBadRequest,
		Message: message,
	})
}
```

- [ ] **Step 5: Commit**

```bash
git add . && git commit -m "feat: add database connection and base model"
```

---

### Task 1.4: Database Migrations

- [ ] **Step 1: Create initial migration**

Create `migrations/001_init.up.sql` (full schema - see design doc section 2)

- [ ] **Step 2: Create rollback migration**

Create `migrations/001_init.down.sql`:
```sql
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS stock_transfers;
DROP TABLE IF EXISTS outbound_items;
DROP TABLE IF EXISTS outbound_orders;
DROP TABLE IF EXISTS inbound_items;
DROP TABLE IF EXISTS inbound_orders;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS suppliers;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS warehouses;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
```

- [ ] **Step 3: Commit**

```bash
git add . && git commit -m "feat: add database migrations"
```

---

## Phase 2: Authentication & Authorization

### Task 2.1: User Model and Repository

- [ ] **Step 1: Create user model**

Create `internal/model/user.go`:
```go
package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"users,alias:u"`
	BaseModel
	Username string `bun:"username,notnull" json:"username"`
	Password string `bun:"password,notnull" json:"-"`
	Nickname string `bun:"nickname" json:"nickname"`
	Email    string `bun:"email" json:"email"`
	Phone    string `bun:"phone" json:"phone"`
	Status   int    `bun:"status,notnull" json:"status"`
}

type Role struct {
	bun.BaseModel `bun:"roles,alias:r"`
	BaseModel
	Name        string `bun:"name,notnull" json:"name"`
	Code        string `bun:"code,notnull" json:"code"`
	Description string `bun:"description" json:"description"`
	Status      int    `bun:"status,notnull" json:"status"`
}

type Permission struct {
	bun.BaseModel `bun:"permissions,alias:p"`
	BaseModel
	Name        string `bun:"name,notnull" json:"name"`
	Code        string `bun:"code,notnull" json:"code"`
	Resource    string `bun:"resource,notnull" json:"resource"`
	Action      string `bun:"action,notnull" json:"action"`
	Description string `bun:"description" json:"description"`
}

type UserRole struct {
	bun.BaseModel `bun:"user_roles,alias:ur"`
	BaseModel
	UserID int64 `bun:"user_id,notnull" json:"user_id"`
	RoleID int64 `bun:"role_id,notnull" json:"role_id"`
}

type RolePermission struct {
	bun.BaseModel `bun:"role_permissions,alias:rp"`
	BaseModel
	RoleID       int64 `bun:"role_id,notnull" json:"role_id"`
	PermissionID int64 `bun:"permission_id,notnull" json:"permission_id"`
}
```

- [ ] **Step 2: Create user repository**

Create `internal/repository/user.go` with CRUD operations and `GetUserPermissions` method.

- [ ] **Step 3: Commit**

```bash
git add . && git commit -m "feat: add user model and repository"
```

---

### Task 2.2: JWT and Password Utilities

- [ ] **Step 1: Create JWT utilities**

Create `internal/pkg/jwt/jwt.go`:
```go
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret []byte
	expire time.Duration
}

func NewJWT(secret string, expire time.Duration) *JWT {
	return &JWT{
		secret: []byte(secret),
		expire: expire,
	}
}

func (j *JWT) GenerateToken(userID int64, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expire)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
```

- [ ] **Step 2: Create password utilities**

Create `internal/pkg/password/password.go`:
```go
package password

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

- [ ] **Step 3: Commit**

```bash
git add . && git commit -m "feat: add JWT and password utilities"
```

---

### Task 2.3: Auth Service, Handler and Middleware

- [ ] **Step 1: Create auth service** - Login, GetProfile, ChangePassword

- [ ] **Step 2: Create auth handler** - HTTP handlers for auth endpoints

- [ ] **Step 3: Create JWT middleware** - Token validation middleware

- [ ] **Step 4: Commit**

---

### Task 2.4: RBAC Middleware

- [ ] **Step 1: Create RBAC middleware** - Load user permissions, RequirePermission function

- [ ] **Step 2: Commit**

---

### Task 2.5: Role and Permission Management

- [ ] **Step 1: Create role repository** - CRUD and AssignPermissions

- [ ] **Step 2: Create permission repository** - CRUD and BatchCreate

- [ ] **Step 3: Create role service and handler**

- [ ] **Step 4: Commit**

---

## Phase 3: Core Business Modules

### Task 3.1: Warehouse Module
- Model, Repository, Service, Handler for warehouses

### Task 3.2: Location Module
- Model, Repository, Service, Handler for locations

### Task 3.3: Product and Category Module
- Model, Repository, Service, Handler for products and categories

### Task 3.4: Inventory Module
- Model, Repository, Service, Handler for inventory

### Task 3.5: Supplier and Customer Module
- Model, Repository, Service, Handler for suppliers and customers

---

## Phase 4: Order Management

### Task 4.1: Inbound Order Module
- Model, Repository, Service, Handler for inbound orders and items

### Task 4.2: Outbound Order Module
- Model, Repository, Service, Handler for outbound orders and items

### Task 4.3: Stock Transfer Module
- Model, Repository, Service, Handler for stock transfers

---

## Phase 5: Audit System

### Task 5.1: Audit Log Model and Repository

### Task 5.2: Audit Service (automatic logging on all CRUD operations)

### Task 5.3: Audit Handler (query endpoints)

---

## Phase 6: Router and Main

### Task 6.1: Router Setup
- Register all routes with proper middleware chain

### Task 6.2: Main Entry
- Load config, initialize DB, services, handlers, start server

---

## Phase 7: Frontend (React + TypeScript)

### Task 7.1: Initialize React Project
- Vite + React + TypeScript + Ant Design setup

### Task 7.2: Layout and Authentication Pages
- Main layout, login page, change password

### Task 7.3: Business Pages
- Warehouse, Location, Product, Inventory, Order pages

### Task 7.4: System Pages
- User, Role, Permission, Audit log pages

---

## Phase 8: Static File Embedding

### Task 8.1: Embed static files into Go binary
- Create embed.go with go:embed directive

### Task 8.2: Serve static files from Gin

---

## Detailed Implementation Notes

For complete code of each task, refer to the design document at:
`docs/superpowers/specs/2026-04-25-warehouse-design.md`

Each task follows the pattern:
1. Write model/repository/service/handler
2. Add appropriate tests
3. Commit with descriptive message
