// ./app/queries/book_query.go

package queries

import (
	"github.com/rvic/adserver-simple/app/models"

	"github.com/jmoiron/sqlx"
)

type CustomerQueries struct {
	*sqlx.DB
}

func (q *CustomerQueries) GetCustomers() ([]models.Customer, error) {
	customers := []models.Customer{}
	query := `SELECT * FROM customers`
	err := q.Select(&customers, query)
	if err != nil {
		//fmt.Println(err)
		return customers, err
	}

	return customers, nil
}

func (q *CustomerQueries) GetCustomer(id string) (models.Customer, error) {
	customer := models.Customer{}
	query := `SELECT * FROM customers WHERE id = $1`

	err := q.Get(&customer, query, id)
	if err != nil {
		return customer, err
	}

	return customer, nil
}

// CreateBook method for creating book by given Book object.
func (q *CustomerQueries) AddCustomer(c *models.Customer) error {
	query := `INSERT INTO customers VALUES ($1, $2, $3)`
	_, err := q.Exec(query, c.ID, c.Name, c.Balance)
	if err != nil {
		return err
	}

	return nil
}

// UpdateBook method for updating book by given Book object.
func (q *CustomerQueries) UpdateCustomer(c *models.CustomerUpd) error {
	query := `UPDATE customers SET balance = $2 WHERE id = $1`

	_, err := q.Exec(query, c.ID, c.Balance)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBook method for delete book by given ID.
func (q *CustomerQueries) DeleteCustomer(id string) error {
	query := `DELETE FROM customers WHERE id = $1`

	_, err := q.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
