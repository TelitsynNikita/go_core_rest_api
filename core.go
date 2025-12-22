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
	DB Database
}

func NewApiService() *ApiService {
	return &ApiService{}
}

func (a *ApiService) InitApiService() {
	var (
		configSignal   = make(chan struct{})
		apisYamlSignal = make(chan struct{})
		dbInitSignal   = make(chan struct{})
		serverSignal   = make(chan struct{})
	)
	go func() {
		a.initConfigs()
		configSignal <- struct{}{}
	}()

	go func() {
		<-configSignal
		a.initApisYaml()
		apisYamlSignal <- struct{}{}
	}()

	go func() {
		<-apisYamlSignal
		a.initDB(a.Config.DBConfig)
		dbInitSignal <- struct{}{}
	}()

	go func() {
		<-dbInitSignal
		a.initServer(&a.DB)
		serverSignal <- struct{}{}
	}()

	<-serverSignal
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
