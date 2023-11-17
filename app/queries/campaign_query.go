package queries

import (
	"github.com/rvic/adserver-simple/app/models"

	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CampaignQueries struct {
	*sqlx.DB
}

func (q *CampaignQueries) GetCampaigns(customer_id string) ([]models.Campaign, error) {
	campaigns := []models.Campaign{}
	campaignRecs := []models.CampaignTableRec{}
	query := `SELECT id, customer_id, creative, views FROM campaigns where customer_id=$1`
	err := q.Select(&campaignRecs, query, customer_id)
	if err != nil {
		//fmt.Println(err)
		return campaigns, err
	}

	for i := 0; i < len(campaignRecs); i++ {
		var campaign = models.Campaign{}
		campaign.ID = campaignRecs[i].ID
		campaign.CustomerID = customer_id
		campaign.Creative = campaignRecs[i].Creative
		campaign.Views = campaignRecs[i].Views

		campaignCountries := []models.CampaignCountry{}
		query_countries := `select cn.name from campaign_countries cc ` +
			`left join countries cn on cn.id=cc.country_id where cc.campaign_id=$1`
		err = q.Select(&campaignCountries, query_countries, campaignRecs[i].ID)
		if err == nil {
			for _, campaignCountry := range campaignCountries {
				campaign.Countries = append(campaign.Countries, campaignCountry.Country)
			}
		}

		campaignDevices := []models.CampaignDevice{}
		query_devices := `select d.name from campaign_devices cd ` +
			`left join devices d on d.id=cd.device_id where cd.campaign_id=$1`
		err = q.Select(&campaignDevices, query_devices, campaignRecs[i].ID)
		if err == nil {
			for _, campaignDevice := range campaignDevices {
				campaign.Devices = append(campaign.Devices, campaignDevice.Device)
			}
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil

}

func (q *CampaignQueries) AddCampaign(c *models.Campaign) error {
	tx, err := q.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO campaigns VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(query, c.ID, c.CustomerID, c.Creative, c.Views)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := 0; i < len(c.Countries); i++ {
		row := tx.QueryRow("SELECT id FROM countries WHERE name=$1", c.Countries[i])
		var country_id string
		err = row.Scan(&country_id)
		if err == sql.ErrNoRows {
			query := `INSERT INTO countries VALUES ($1, $2)`
			country_id = uuid.New().String()
			_, err := tx.Exec(query, country_id, c.Countries[i])
			if err != nil {
				tx.Rollback()
				return err
			}
		} else if err != nil {
			tx.Rollback()
			return err
		}
		query := `INSERT INTO campaign_countries VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, uuid.New().String(), c.ID, country_id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for i := 0; i < len(c.Devices); i++ {
		row := tx.QueryRow("SELECT id FROM devices WHERE name=$1", c.Devices[i])
		var device_id string
		err = row.Scan(&device_id)
		if err == sql.ErrNoRows {
			query := `INSERT INTO devices VALUES ($1, $2)`
			device_id = uuid.New().String()
			_, err := tx.Exec(query, device_id, c.Devices[i])
			if err != nil {
				tx.Rollback()
				return err
			}
		} else if err != nil {
			tx.Rollback()
			return err
		}

		query = `INSERT INTO campaign_devices VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, uuid.New().String(), c.ID, device_id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
