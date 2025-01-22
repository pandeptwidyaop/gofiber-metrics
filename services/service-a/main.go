package main

import (
	"context"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"golang-metrics/services"
	"golang-metrics/stuff"
)

// Service A to emulated api gateway

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := services.InitTracing(ctx, "service-a")
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Service-A",
	})
	app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	})
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(otelfiber.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		stuff.DoSomeWork(c.UserContext())
		return c.SendString("Hello World")
	})

	log.Fatal(app.Listen(":8001"))
}
