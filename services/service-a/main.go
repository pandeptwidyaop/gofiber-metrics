package main

import (
	"context"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"golang-metrics/services"
	"golang-metrics/stuff"
	"net/http"
)

// Service A to emulated api gateway

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/**
	Initiate the tracing
	*/
	_, err := services.InitTracing(ctx, "service-a")
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Service-A",
	})

	/**
	Adding request ID

	*/
	app.Use(requestid.New())

	/**
	Use the same context fir fiber.Ctx
	*/
	app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	})

	app.Use(logger.New())
	/**
	IMPORTANT! use otelfiber middleware
	*/
	app.Use(otelfiber.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		stuff.DoSomeWork(c.UserContext())

		httpClient := &http.Client{}

		req, err := http.NewRequest("GET", "http://localhost:8002/", nil)
		if err != nil {
			return err
		}

		otel.GetTextMapPropagator().Inject(c.UserContext(), propagation.HeaderCarrier(req.Header))

		httpResp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		defer httpResp.Body.Close()

		return c.SendString("Hello World")
	})

	log.Fatal(app.Listen(":8001"))
}
