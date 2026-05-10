package go_core_rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"db_name"`
	Driver   string `json:"driver"`
	SSLMode  string `json:"ssl_mode"`

	// Pool settings; zero values pick sensible defaults (see applyPoolSettings).
	MaxOpenConns       int `json:"max_open_conns"`
	MaxIdleConns       int `json:"max_idle_conns"`
	ConnMaxLifetimeSec int `json:"conn_max_lifetime_sec"`
	ConnMaxIdleTimeSec int `json:"conn_max_idle_time_sec"`
}

type Database struct {
	SQLX *sqlx.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func NewDBConnection(config DBConfig) (*sqlx.DB, error) {
	dbConnect, err := sqlx.Open(config.Driver, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", config.Host, config.Port, config.User, config.DBName, config.Password, config.SSLMode))
	if err != nil {
		return nil, err
	}

	if err = dbConnect.Ping(); err != nil {
		return nil, err
	}

	applyPoolSettings(dbConnect, config)

	return dbConnect, nil
}

func applyPoolSettings(db *sqlx.DB, cfg DBConfig) {
	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 25
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 5
	}
	if maxIdle > maxOpen {
		maxIdle = maxOpen
	}

	lifetimeSec := cfg.ConnMaxLifetimeSec
	if lifetimeSec <= 0 {
		lifetimeSec = 300
	}
	idleSec := cfg.ConnMaxIdleTimeSec
	if idleSec <= 0 {
		idleSec = 900
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(lifetimeSec) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(idleSec) * time.Second)
}

func (db *Database) SelectFunction(ctx context.Context, functionName string, body []byte) (interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var result []byte
	var err error

	if body == nil {
		query := fmt.Sprintf("SELECT %s()", functionName)
		err = db.SQLX.QueryRowContext(ctx, query).Scan(&result)
	} else {
		query := fmt.Sprintf("SELECT %s($1)", functionName)
		err = db.SQLX.QueryRowContext(ctx, query, string(body)).Scan(&result)
	}
	if err != nil {
		return nil, err
	}

	var bodyJson interface{}
	err = json.Unmarshal(result, &bodyJson)
	if err != nil {
		return nil, err
	}

	return bodyJson, nil
}
