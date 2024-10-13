package http

import (
	"checkout-service/internal/infrastructure/container"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func HTTPRouteInit(containerConf *container.Container) {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &containerConf.Logger.Log,
		Fields: []string{"locals:requestid", "method", "path", "pid", "status", "resBody",
			"latency", "reqHeaders", "body"},
		WrapHeaders: true,
		SkipURIs:    []string{"/"},
	}))

	SetupRouter(app, *containerConf)

	// Start server
	port := fmt.Sprintf("%s:%d", containerConf.Apps.Host, containerConf.Apps.HttpPort)
	if containerConf.Apps.HttpPort == 0 {
		port = ":8000"
	}
	go func() {
		if err := app.Listen(port); err != nil && err != http.ErrServerClosed {
			app.Server().Logger.Printf("shutting down the server : ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		app.Server().Logger.Printf("shutting down the server :", err)
	}

}
