package product

import "encore.dev/storage/sqldb"

var db = sqldb.NewDatabase("product", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
