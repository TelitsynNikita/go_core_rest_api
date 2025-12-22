package go_core_rest_api

import (
	"encoding/json"
	"fmt"

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
}

type Database struct {
	sqlx *sqlx.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func NewDBConnection(config DBConfig) (*sqlx.DB, error) {
	dbConnect, err := sqlx.Open(config.Driver, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", config.Host, config.Port, config.User, config.DBName, config.Password, config.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := dbConnect.Ping(); err != nil {
		return nil, err
	}

	return dbConnect, nil
}

func (db *Database) SelectFunction(functionName string, body []byte) (interface{}, error) {
	tx, err := db.sqlx.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var result []byte
	var query string

	if body == nil {
		query = fmt.Sprintf("SELECT %s()", functionName)
		row := db.sqlx.QueryRow(query)
		if row.Err() != nil {
			return nil, err
		}

		err = row.Scan(&result)
		if err != nil {
			return nil, err
		}
	} else {
		query = fmt.Sprintf("SELECT %s($1)", functionName)
		row := db.sqlx.QueryRow(query, string(body))
		if row.Err() != nil {
			return nil, err
		}

		err = row.Scan(&result)
		if err != nil {
			return nil, err
		}
	}

	var something interface{}
	err = json.Unmarshal(result, &something)
	if err != nil {
		return nil, err
	}

	return something, nil
}
