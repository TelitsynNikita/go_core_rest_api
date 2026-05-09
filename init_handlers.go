package go_core_rest_api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

func InitServerApp(groups map[string]map[string]ApiSettings, config Config, db *Database) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	for groupName, group := range groups {
		apiGroup := app.Group(groupName)

		for url, api := range group {
			if strings.ToLower(group[url].Type) == "custom" {

			} else {
				makeHandler(url, api, apiGroup, db)
			}
		}
	}

	return app
}

func makeHandler(url string, apiSetting ApiSettings, apiGroup fiber.Router, db *Database) {
	switch strings.ToLower(apiSetting.Method) {
	case strings.ToLower(fiber.MethodGet):
		apiGroup.Get(url, func(c *fiber.Ctx) error {
			result, err := db.SelectFunction(apiSetting.Call, nil)
			if err != nil {
				logrus.Error(err)
				return err
			}

			return c.JSON(result)
		})
	case strings.ToLower(fiber.MethodPost):
		apiGroup.Post(url, func(c *fiber.Ctx) error {
			body, err := parseBody(apiSetting.Body, c.Body(), apiSetting.IsSlice)
			if err != nil {
				logrus.Error(err)
				return err
			}

			result, err := db.SelectFunction(apiSetting.Call, body)
			if err != nil {
				logrus.Error(err)
				return err
			}

			return c.JSON(result)
		})
	case strings.ToLower(fiber.MethodPut):
		apiGroup.Put(url, func(c *fiber.Ctx) error {
			body, err := parseBody(apiSetting.Body, c.Body(), apiSetting.IsSlice)
			if err != nil {
				logrus.Error(err)
				return err
			}

			result, err := db.SelectFunction(apiSetting.Call, body)
			if err != nil {
				logrus.Error(err)
				return err
			}

			return c.JSON(result)
		})
	case strings.ToLower(fiber.MethodDelete):
		apiGroup.Delete(url, func(c *fiber.Ctx) error {
			body, err := parseBody(apiSetting.Body, c.Body(), apiSetting.IsSlice)
			if err != nil {
				logrus.Error(err)
				return err
			}

			result, err := db.SelectFunction(apiSetting.Call, body)
			if err != nil {
				logrus.Error(err)
				return err
			}

			return c.JSON(result)
		})
	}
}

func parseBody(bodySettings map[string]map[string]string, bodyBytes []byte, isSlice bool) ([]byte, error) {
	if len(bodySettings) == 0 {
		return bodyBytes, nil
	}

	if len(bodyBytes) == 0 {
		return nil, nil
	}

	var fieldsOfCustomStruct []reflect.StructField
	for field, fieldSettingsMap := range bodySettings {
		tag, ok := fieldSettingsMap["tag"]
		if !ok {
			return nil, fmt.Errorf("there's no tag by field: %s", field)
		}

		typeOfField, ok := fieldSettingsMap["type"]
		if !ok {
			return nil, fmt.Errorf("there's no type by field: %s", field)
		}

		var reflectTypeOfField reflect.Type

		if fieldSettingsMap["is_slice"] == "true" {
			switch typeOfField {
			case "string":
				reflectTypeOfField = reflect.TypeOf([]string{})
			case "int":
				reflectTypeOfField = reflect.TypeOf([]int{})
			case "bool":
				reflectTypeOfField = reflect.TypeOf([]bool{})
			case "float64":
				reflectTypeOfField = reflect.TypeOf([]float64{})
			case "float32":
				reflectTypeOfField = reflect.TypeOf([]float32{})
			case "int64":
				reflectTypeOfField = reflect.TypeOf([]int64{})
			}
		} else {
			switch typeOfField {
			case "string":
				reflectTypeOfField = reflect.TypeOf("")
			case "int":
				reflectTypeOfField = reflect.TypeOf(0)
			case "bool":
				reflectTypeOfField = reflect.TypeOf(false)
			case "float64":
				reflectTypeOfField = reflect.TypeOf(float64(0))
			case "float32":
				reflectTypeOfField = reflect.TypeOf(float32(0))
			case "int64":
				reflectTypeOfField = reflect.TypeOf(int64(0))
			}
		}

		var fieldOfCustomStruct = reflect.StructField{
			Name: strings.ToUpper(field[:1]) + field[1:],
			Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\" db:\"%s\" validate:\"%s\"", field, field, tag)),
			Type: reflectTypeOfField,
		}

		fieldsOfCustomStruct = append(fieldsOfCustomStruct, fieldOfCustomStruct)
	}

	if isSlice {
		var unmarshalledBody []map[string]interface{}
		err := jsoniter.Unmarshal(bodyBytes, &unmarshalledBody)
		if err != nil {
			return nil, err
		}
		// Создаем новый тип слайса структур
		extendedSliceType := reflect.SliceOf(reflect.StructOf(fieldsOfCustomStruct))

		// Создаем экземпляр нового слайса
		extendedSliceValue := reflect.New(extendedSliceType).Elem()

		for _, item := range unmarshalledBody {
			for key, value := range item {
				// Создаем новый тип структуры
				extendedType := reflect.StructOf(fieldsOfCustomStruct)

				// Создаем экземпляр новой структуры
				extendedValue := reflect.New(extendedType).Elem()

				fieldOfCustomStructValue := extendedValue.FieldByName(strings.ToUpper(key[:1]) + key[1:])
				if fieldOfCustomStructValue.IsValid() {
					if fieldOfCustomStructValue.Type() == reflect.TypeOf(value) {
						fieldOfCustomStructValue.Set(reflect.ValueOf(value))
					} else {
						return nil, fmt.Errorf("wrong type for field %s value: %v", key, value)
					}
				} else {
					return nil, fmt.Errorf("there's no field: %s", key)
				}

				validate := validator.New(validator.WithRequiredStructEnabled())
				if err = validate.Struct(extendedValue.Interface()); err != nil {
					return nil, err
				}

				extendedSliceValue = reflect.Append(extendedSliceValue, extendedValue)
			}
		}

		bodyBytes, err = json.Marshal(extendedSliceValue.Interface())
		if err != nil {
			return nil, err
		}

		return bodyBytes, nil
	} else {
		var unmarshalledBody map[string]interface{}
		err := jsoniter.Unmarshal(bodyBytes, &unmarshalledBody)
		if err != nil {
			return nil, err
		}

		// Создаем новый тип структуры
		extendedType := reflect.StructOf(fieldsOfCustomStruct)

		// Создаем экземпляр новой структуры
		extendedValue := reflect.New(extendedType).Elem()

		fmt.Println(extendedValue)

		for key, value := range unmarshalledBody {
			fieldOfCustomStructValue := extendedValue.FieldByName(strings.ToUpper(key[:1]) + key[1:])
			if fieldOfCustomStructValue.IsValid() {
				if fieldOfCustomStructValue.Type() == reflect.TypeOf(value) {
					fieldOfCustomStructValue.Set(reflect.ValueOf(value))
				} else {
					switch fieldOfCustomStructValue.Kind() {
					case reflect.String:
						fieldOfCustomStructValue.Set(reflect.ValueOf(fmt.Sprintf("%v", value)))
					case reflect.Int:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int(value.(float32))))
						}
					case reflect.Int64:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int64(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int64(value.(float32))))
						}
					case reflect.Int32:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int32(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int32(value.(float32))))
						}
					case reflect.Int16:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int16(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int16(value.(float32))))
						}
					case reflect.Int8:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int8(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(int8(value.(float32))))
						}
					case reflect.Uint:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint(value.(float32))))
						}
					case reflect.Uint64:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint64(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint64(value.(float32))))
						}
					case reflect.Uint32:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint32(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint32(value.(float32))))
						}
					case reflect.Uint16:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint16(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint16(value.(float32))))
						}
					case reflect.Uint8:
						if reflect.TypeOf(value).Kind() == reflect.Float64 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint8(value.(float64))))
						} else if reflect.TypeOf(value).Kind() == reflect.Float32 {
							fieldOfCustomStructValue.Set(reflect.ValueOf(uint8(value.(float32))))
						}
					default:
						return nil, fmt.Errorf("wrong type for field %s value: %v", key, value)
					}
				}
			} else {
				return nil, fmt.Errorf("there's no field: %s", key)
			}
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		if err = validate.Struct(extendedValue.Interface()); err != nil {
			return nil, err
		}

		extendedBodyBytes, err := json.Marshal(extendedValue.Interface())
		if err != nil {
			return nil, err
		}

		return extendedBodyBytes, nil
	}
}
