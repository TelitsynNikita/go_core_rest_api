package go_core_rest_api

import (
	"github.com/gofiber/fiber/v2"
)

type ApiService struct {
	YamlConfig map[string]map[string]ApiSettings
	Server     *fiber.App
}

func NewApiService() *ApiService {
	return &ApiService{}
}

func (a *ApiService) InitApiService() {
	a.initApisYaml()
	a.initServer()
	a.initHandlers()
}

func (a *ApiService) RunService() error {
	return a.Server.Listen(":8080")
}

func (a *ApiService) initApisYaml() {
	yamlConfig, err := ReadApisYaml()
	if err != nil {
		panic(err)
	}

	a.YamlConfig = yamlConfig
}

func (a *ApiService) initServer() {
	a.Server = fiber.New()
}

func (a *ApiService) initHandlers() {
	a.Server.Server().Handler = InitHandlers(a.YamlConfig)
}
