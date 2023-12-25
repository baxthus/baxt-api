package endpoint

import "github.com/gofiber/fiber/v2"

type Endpoint struct {
	Method  string
	Url     string
	Handler func(ctx *fiber.Ctx) error
}
