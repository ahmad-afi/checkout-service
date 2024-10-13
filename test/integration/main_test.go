package integration_test

import (
	delivery "checkout-service/internal/delivery/http"
	"checkout-service/internal/infrastructure/container"
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

func TestMain(m *testing.M) {
	cont := container.NewContainer()

	app = fiber.New(fiber.Config{
		// Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			ctx.Status(fiber.StatusInternalServerError)
			return ctx.SendString("Error : " + err.Error())
		},
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &cont.Logger.Log,
		Fields: []string{"locals:requestid", "method", "path", "pid", "status", "resBody",
			"latency", "reqHeaders", "body"},
		WrapHeaders: true,
		SkipURIs:    []string{"/"},
	}))

	app.Get("", func(ctx *fiber.Ctx) error {
		fmt.Println("c.Err()", ctx.Context().Err())
		if err := ctx.Context().Err(); err != nil {
			if err == context.Canceled {
				fmt.Println("Context was canceled")
				return ctx.Status(fiber.StatusInternalServerError).SendString("Context was canceled")
			} else if err == context.DeadlineExceeded {
				fmt.Println("Context deadline exceeded")
			}
		} else {
			fmt.Println("Context still active")
		}
		return ctx.Status(fiber.StatusOK).SendString("Server is up and running")
	})

	delivery.SetupRouter(app, *cont)
	m.Run()
}

func TestHealth(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request, int(600*time.Second))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	fmt.Println("bytes", string(bytes))
	// assert.Equal(t, "Hello World", string(bytes))
}
