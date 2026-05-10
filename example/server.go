package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TelitsynNikita/go_core_rest_api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	apiService := go_core_rest_api.NewApiService()

	apiService.InitApiService()

	apiService.InitCustomHandler("api", "custom", fiber.MethodPost, func(c *fiber.Ctx) error {
		var result []byte

		query := fmt.Sprintf("SELECT %s()", "equipment.get_all_cs_equipment")
		err := apiService.DB.SQLX.QueryRowContext(context.Background(), query).Scan(&result)
		if err != nil {
			return err
		}

		var bodyJson interface{}
		err = json.Unmarshal(result, &bodyJson)
		if err != nil {
			return err
		}

		return c.JSON(bodyJson)
	})

	if err := apiService.RunService(); err != nil {
		fmt.Println(err)
	}
}
