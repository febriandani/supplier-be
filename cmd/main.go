package main

import (
	"log"

	"supplier-be/cmd/routes"
	handler "supplier-be/handler"
	infra "supplier-be/infra"
	repository "supplier-be/repository"
	service "supplier-be/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

const (
	dbDriver = "postgres"
	dbSource = "host=localhost port=5432 user=postgres password=junior34 dbname=suppliers_db sslmode=disable connect_timeout=3"
	port     = ":8080" // Replace with the desired port number
)

func main() {
	logger := infra.NewLogger()

	// Connect to the database
	db, err := sqlx.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Initialize the repository, service, and handler
	supplierRepo := repository.NewSupplierRepository(db)
	supplierService := service.NewSupplierService(supplierRepo, db, logger)
	supplierHandler := handler.NewSupplierHandler(supplierService, logger)

	// Create a new Fiber instance
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app, supplierHandler, logger)

	// Start the Fiber server
	logger.WithField("StartApp", "gofiber").Info("server listen to port ", port)
	logger.Fatal(app.Listen(port))
	logger.WithField("StartApp", "gofiber").Info("server listen to port ", port)

}
