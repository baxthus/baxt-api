package main

import (
	"baxt-api/endpoints"
	"baxt-api/internal/endpoint"
	"baxt-api/internal/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var app = fiber.New()

func init() {
	for _, route := range routes {
		app.Add(route.Method, route.Url, route.Handler)
	}
}

func main() {
	log.Fatal(app.Listen(":" + env.Port))
}

var routes = []endpoint.Endpoint{
	endpoints.IpEndpoint,
}
