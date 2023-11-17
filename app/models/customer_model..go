package models

type Customer struct {
	ID      string `db:"id" json:"id" validate:"required,uuid"`
	Balance int    `db:"balance" json:"balance" validate:"required"`
	Name    string `db:"name" json:"name" validate:"required,lte=255"`
}

type CustomerCrt struct {
	Balance int    `json:"balance" validate:"required"`
	Name    string `json:"name" validate:"required,lte=255"`
}

type CustomerUpd struct {
	ID      string `json:"id" validate:"required,uuid"`
	Balance int    `json:"balance" validate:"required"`
}

type CustomerDel struct {
	ID string `json:"id" validate:"required,uuid"`
}
