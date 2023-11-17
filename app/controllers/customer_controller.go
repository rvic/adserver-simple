// ./app/controllers/customer_controller.go

package controllers

import (
	"time"

	"github.com/rvic/adserver-simple/app/models"
	"github.com/rvic/adserver-simple/pkg/utils"
	"github.com/rvic/adserver-simple/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCustomers func gets all exising customers.
// @Description Get all existing customers.
// @Summary get all existsing customers
// @Tags Customers
// @Accept json
// @Produce json
// @Success 200 {array} models.Customer
// @Router /api/v1/customers [get]
func GetCustomers(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	customers, err := db.GetCustomers()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":     true,
			"msg":       "customers were not found",
			"count":     0,
			"customers": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":     false,
		"msg":       nil,
		"count":     len(customers),
		"customers": customers,
	})
}

// GetCustomer func gets customer by given ID or 404 error.
// @Description Get customer by given ID.
// @Summary get customer by given ID
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Customer
// @Router /api/v1/customer/{id} [get]
func GetCustomer(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	customer, err := db.GetCustomer(id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "customer with the given ID is not found",
			"customer": nil,
		})
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"customer": customer,
	})
}

// AddCustomer func for creates a new customer.
// @Description Create a new customer.
// @Summary create a new customer
// @Tags Customer
// @Accept json
// @Produce json
// @param customer body models.CustomerCrt true "Customer"
// @Success 200 {object} models.Customer
// @Security ApiKeyAuth
// @Router /api/v1/customer [post]
func AddCustomer(c *fiber.Ctx) error {

	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	expires := claims.Expires

	// Checking, if now time greater than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	customerCrt := &models.CustomerCrt{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(customerCrt); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	customer := &models.Customer{}
	customer.Name = customerCrt.Name
	customer.Balance = customerCrt.Balance

	validate := utils.NewValidator()

	customer.ID = uuid.New().String()

	if err := validate.Struct(customer); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := db.AddCustomer(customer); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"customer": customer,
	})
}

// UpdateCustomer func for updates customer by given ID.
// @Description Update customer.
// @Summary update customer
// @Tags Customer
// @Accept json
// @Produce json
// @param customer body models.CustomerUpd true "Customer"
// @Success 201 {string} status "ok"
// @Security ApiKeyAuth
// @Router /api/v1/customer [put]
func UpdateCustomer(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	customerUpdate := &models.CustomerUpd{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(customerUpdate); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	_, err = db.GetCustomer(customerUpdate.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "customer with this ID not found",
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(customerUpdate); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := db.UpdateCustomer(customerUpdate); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 201.
	return c.SendStatus(fiber.StatusCreated)
}

// DeleteCustomer func for deletes customer by given ID.
// @Description Delete customer by given ID.
// @Summary delete customer by given ID
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body models.CustomerDel true "Customer"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /api/v1/customer [delete]
func DeleteCustomer(c *fiber.Ctx) error {

	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	customerDelete := &models.CustomerDel{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(customerDelete); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	_, err = db.GetCustomer(customerDelete.ID)
	if err != nil {
		// Return status 404 and customer not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "customer with this ID not found",
		})
	}

	if err := db.DeleteCustomer(customerDelete.ID); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
