package main

import (
	"baxt-api/endpoints"
	"baxt-api/internal/endpoint"
	"baxt-api/internal/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var app = fiber.New()

func init() {
	for _, route := range routes {
		app.Add(route.Method, route.Url, route.Handler)
	}
}

func main() {
	app.Use(
		cors.New(),
		compress.New(),
		recover.New(),
	)

	app.Static("/", "/public")

	log.Fatal(app.Listen(":" + env.Port))
}

var routes = []endpoint.Endpoint{
	endpoints.IpEndpoint,
}
