package go_core_rest_api

import (
	"fmt"
	"strings"

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

func (a *ApiService) InitCustomHandler(group, url, method string, fn func(c *fiber.Ctx) error) {
	switch strings.ToUpper(method) {
	case fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodPut,
		fiber.MethodDelete:
	default:
		panic("Unsupported http method: " + method)
	}

	a.Server.Add(strings.ToUpper(method), fmt.Sprintf("%s/%s", group, url), fn)
}

func (a *ApiService) InitApiService() {
	a.initConfigs()
	a.initApisYaml()
	a.initDB()
	a.initServer(&a.DB)
}

func (a *ApiService) RunService() error {
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

func (a *ApiService) initDB() {
	db := NewDatabase()
	sqlx, err := NewDBConnection(a.Config.DBConfig)
	if err != nil {
		panic(err)
	}

	db.SQLX = sqlx
	a.DB = *db
}

func (a *ApiService) initServer(db *Database) {
	a.Server = InitServerApp(a.YamlConfig, a.Config, db)
}
