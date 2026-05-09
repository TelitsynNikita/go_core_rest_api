package go_core_rest_api

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ApiSettings struct {
	Body    map[string]map[string]string `yaml:"body"`
	IsSlice bool                         `yaml:"is_slice"`
	Type    string                       `yaml:"type"`
	Method  string                       `yaml:"method"`
	Call    string                       `yaml:"call"`
}

func ReadApisYaml() (map[string]map[string]ApiSettings, error) {
	apisYaml, err := os.ReadFile("apis.yaml")
	if err != nil {
		return nil, err
	}

	var apis map[string]map[string]ApiSettings
	err = yaml.Unmarshal(apisYaml, &apis)
	if err != nil {
		return nil, err
	}

	err = validateApis(apis)
	if err != nil {
		return nil, err
	}

	return apis, nil
}

func validateApis(groups map[string]map[string]ApiSettings) error {
	for _, group := range groups {
		for key, api := range group {
			if api.Method == "" {
				return errors.New(fmt.Sprintf("There's no method in api: %s. Method is required in apis.yaml", key))
			}

			if strings.ToLower(api.Type) != "custom" && strings.ToLower(api.Type) != "postgresql" {
				return errors.New(fmt.Sprintf("Unsupported type %s", api.Type))
			}
		}
	}

	return nil
}
