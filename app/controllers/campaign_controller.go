package controllers

import (
	"github.com/rvic/adserver-simple/app/models"
	"github.com/rvic/adserver-simple/pkg/utils"
	"github.com/rvic/adserver-simple/platform/database"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCampaigns func gets campaigns for certain customer.
// @Description Get customer campaigns.
// @Summary  Gets campaigns for customer
// @Tags Campaigns
// @Accept json
// @Produce json
// @Param customer_id query string true "Customer ID"
// @Success 200 {array} models.Campaign
// @Router /api/v1/campaigns [get]
func GetCampaigns(c *fiber.Ctx) error {
	customer_id, err := uuid.Parse(c.Query("customer_id"))
	fmt.Println(err)
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	campaigns, err := db.GetCampaigns(customer_id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":     true,
			"msg":       "campaigns were not found",
			"count":     0,
			"campaigns": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":     false,
		"msg":       nil,
		"count":     len(campaigns),
		"campaigns": campaigns,
	})
}

// AddCampaign func for creates a new campaign.
// @Description Create a new campaign.
// @Summary create a new campaign
// @Tags Campaign
// @Accept json
// @Produce json
// @param campaign body models.CampaignCrt true "Campaign"
// @Success 200 {object} models.Campaign
// @Security ApiKeyAuth
// @Router /api/v1/campaign [post]
func AddCampaign(c *fiber.Ctx) error {

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

	campaignCrt := &models.CampaignCrt{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(campaignCrt); err != nil {
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

	campaign := &models.Campaign{}

	campaign.ID = uuid.New().String()
	campaign.CustomerID = campaignCrt.CustomerID
	campaign.Creative = campaignCrt.Creative
	campaign.Views = campaignCrt.Views
	campaign.Countries = campaignCrt.Countries
	campaign.Devices = campaignCrt.Devices

	validate := utils.NewValidator()

	if err := validate.Struct(campaign); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := db.AddCampaign(campaign); err != nil {
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
		"campaign": campaign,
	})
}
