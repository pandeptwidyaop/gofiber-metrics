package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang-metrics/metrics"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		<-sig
		log.Info("receive signal, shutting down")
		cancel()
	}()

	// Start App Server
	go func() {
		app := fiber.New(
			fiber.Config{
				Prefork: false,
			})

		app.Use(logger.New())

		app.Use(func(ctx *fiber.Ctx) error {
			start := time.Now()
			err := ctx.Next()
			duration := time.Since(start)
			metrics.RequestDuration.With(map[string]string{
				"path": ctx.Route().Path,
			}).Observe(duration.Seconds())
			metrics.RequestTotal.With(map[string]string{
				"path": ctx.Route().Path,
				"code": strconv.Itoa(ctx.Response().StatusCode()),
			}).Inc()

			return err
		})

		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Hello, World!",
			})
		})

		app.Post("/post", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Hello, World!",
			})
		})

		log.Fatal(app.Listen(":3000"))
	}()

	// Start Metrics Server
	go func() {
		app := fiber.New()

		app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

		panic(app.Listen(":3001"))
	}()

	<-ctx.Done()
	log.Info("receive signal, shutting down")

}
