package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	handler "supplier-be/handler"
)

// SetupRoutes sets up the routes for the application.
func SetupRoutes(app *fiber.App, supplierHandler *handler.SupplierHandler, log *logrus.Logger) {
	app.Post("/supplier", supplierHandler.CreateSupplierHandler)
	app.Post("/suppliers", supplierHandler.GetListSupplier)
	app.Get("/", func(c *fiber.Ctx) error {
		response := map[string]string{
			"message": "Health Check OK",
		}
		return c.JSON(response)
	})

}
