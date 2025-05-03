package app

import (
	"go-image-compression/internal/config"
	"go-image-compression/internal/controller/v1/http"

	"github.com/gofiber/fiber/v2"
)

func MustRun() error {
	// ctx, cancel := context.WithCancel(context.Background()) Temporary commented
	// defer cancel()

	cfg, err := config.MustLoad()
	if err != nil {
		return err
	}

	handler := http.NewHandler()

	app := fiber.New()
	handler.SetupRoutes(app)

	if err := app.Listen(cfg.HTTP.Port); err != nil {
		return err
	}

	return nil
}
