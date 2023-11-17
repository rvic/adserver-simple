// ./pkg/routes/private_routes.go

package routes

import (
	"github.com/rvic/adserver-simple/app/controllers"
	"github.com/rvic/adserver-simple/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/customer", middleware.JWTProtected(), controllers.AddCustomer)

	// Routes for PUT method:
	route.Put("/customer", middleware.JWTProtected(), controllers.UpdateCustomer)

	// Routes for DELETE method:
	route.Delete("/customer", middleware.JWTProtected(), controllers.DeleteCustomer)

	route.Post("/campaign", middleware.JWTProtected(), controllers.AddCampaign)
}
