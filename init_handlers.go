package go_core_rest_api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func InitHandlers(groups map[string]map[string]ApiSettings) fasthttp.RequestHandler {
	app := fiber.New()
	for groupName, group := range groups {
		apiGroup := app.Group(groupName)

		for url, api := range group {
			if strings.ToLower(api.Method) == strings.ToLower(fiber.MethodGet) {
				apiGroup.Get(url, func(c *fiber.Ctx) error {
					return c.JSON(fiber.Map{
						"api": api,
					})
				})
			} else if strings.ToLower(api.Method) == strings.ToLower(fiber.MethodPost) {
				apiGroup.Post(url, func(c *fiber.Ctx) error {
					return c.JSON(fiber.Map{
						"api": api,
					})
				})
			}
		}
	}

	return app.Handler()
}
