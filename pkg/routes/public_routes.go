// ./pkg/routes/private_routes.go

package routes

import (
	"github.com/rvic/adserver-simple/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	route.Get("/customers", controllers.GetCustomers)
	route.Get("/customer/:id", controllers.GetCustomer)
	route.Get("/token/new", controllers.GetNewAccessToken) // create a new access tokens

	route.Get("/campaigns", controllers.GetCampaigns)
}
