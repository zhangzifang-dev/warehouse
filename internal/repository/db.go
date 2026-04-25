package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	"warehouse/internal/config"
)

func BuildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}

func NewDB(cfg *config.Config) (*bun.DB, error) {
	dsn := BuildDSN(cfg.Database)

	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqldb.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	db := bun.NewDB(sqldb, mysqldialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewDBFromDSN(dsn string, maxOpenConns, maxIdleConns int) (*bun.DB, error) {
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxIdleConns)

	db := bun.NewDB(sqldb, mysqldialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
