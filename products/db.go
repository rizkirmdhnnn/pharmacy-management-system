package products

import "encore.dev/storage/sqldb"

// Create the products database and assign it to the "db" variable
var db = sqldb.NewDatabase("products", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
