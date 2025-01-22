package main

import (
	"context"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang-metrics/services"
	"golang-metrics/stuff"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := services.InitTracing(ctx, "service-c")
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Use(logger.New())

	app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)

		return c.Next()
	})

	app.Use(otelfiber.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {
		stuff.DoSomeWork(c.UserContext())

		return c.SendString("Hello, World! From service 3")
	})

	panic(app.Listen(":8003"))

}
