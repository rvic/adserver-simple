package models

type Campaign struct {
	ID         string   `json:"id" validate:"required,uuid"`
	CustomerID string   `json:"customer_id" validate:"required,uuid"`
	Creative   string   `json:"creative" validate:"required,lte=255"`
	Countries  []string `json:"countries" validate:"required"`
	Devices    []string `json:"devices" validate:"required"`
	Views      int      `json:"views" validate:"required"`
}

type CampaignCrt struct {
	CustomerID string   `json:"customer_id" validate:"required,uuid"`
	Creative   string   `json:"creative" validate:"required,lte=255"`
	Countries  []string `json:"countries" validate:"required"`
	Devices    []string `json:"devices" validate:"required"`
	Views      int      `json:"views" validate:"required"`
}

type CampaignTableRec struct {
	ID         string `db:"id"`
	CustomerID string `db:"customer_id"`
	Creative   string `db:"creative"`
	Views      int    `db:"views"`
}

type CampaignCountry struct {
	Country string `db:"name"`
}

type CampaignDevice struct {
	Device string `db:"name"`
}
