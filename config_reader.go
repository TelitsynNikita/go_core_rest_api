package go_core_rest_api

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerConfig struct {
		Port string `json:"port"`
	} `json:"server_config"`
	DBConfig     DBConfig `json:"db_config"`
	SSLMode      bool     `json:"ssl_mode"`
	JWTTokenSalt string   `json:"jwt_token_salt"`
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) ReadConfig() error {
	config, err := os.ReadFile("./configs/config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &c)
	if err != nil {
		return err
	}

	return nil
}
