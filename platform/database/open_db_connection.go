// ./platform/database/open_db_connection.go

package database

import "github.com/rvic/adserver-simple/app/queries"

// Queries struct for collect all app queries.
type Queries struct {
	*queries.CustomerQueries
	*queries.CampaignQueries
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define a new PostgreSQL connection.
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		// Set queries from models:
		CustomerQueries: &queries.CustomerQueries{DB: db},
		CampaignQueries: &queries.CampaignQueries{DB: db},
	}, nil
}
