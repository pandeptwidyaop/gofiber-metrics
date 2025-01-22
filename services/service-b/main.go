package main

import (
	"context"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"golang-metrics/services"
	"golang-metrics/stuff"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := services.InitTracing(ctx, "service-b")
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "service-b",
	})

	app.Use(logger.New())

	app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	})

	app.Use(otelfiber.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {

		client := &http.Client{}
		request, err := http.NewRequest("GET", "http://localhost:8003/", nil)
		if err != nil {
			return err
		}

		otel.GetTextMapPropagator().Inject(c.UserContext(), propagation.HeaderCarrier(request.Header))

		response, err := client.Do(request)
		if err != nil {
			return err
		}

		defer response.Body.Close()

		stuff.DoSomeWork(c.UserContext())

		return nil
	})

	panic(app.Listen(":8002"))
}
