package go_core_rest_api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type ApiService struct {
	YamlConfig map[string]map[string]ApiSettings
	Server     *fiber.App
	Config
	DB             Database
	CustomHandlers map[string]fiber.Handler
}

func NewApiService() *ApiService {
	return &ApiService{}
}

func (a *ApiService) InitCustomHandler(key string, fn func(db *Database)) {

}

func (a *ApiService) InitApiService() {
	a.initConfigs()
	a.initApisYaml()
	a.initDB(a.Config.DBConfig)
}

func (a *ApiService) RunService() error {
	a.initServer(&a.DB)
	return a.Server.Listen(fmt.Sprintf(":%s", a.Config.ServerConfig.Port))
}

func (a *ApiService) initConfigs() {
	a.Config = NewConfig()

	err := a.Config.ReadConfig()
	if err != nil {
		panic(err)
	}
}

func (a *ApiService) initApisYaml() {
	yamlConfig, err := ReadApisYaml()
	if err != nil {
		panic(err)
	}

	a.YamlConfig = yamlConfig
}

func (a *ApiService) initDB(config DBConfig) {
	db := NewDatabase()
	sqlx, err := NewDBConnection(config)
	if err != nil {
		panic(err)
	}

	db.sqlx = sqlx
	a.DB = *db
}

func (a *ApiService) initServer(db *Database) {
	a.Server = InitServerApp(a.YamlConfig, a.Config, db)
}
