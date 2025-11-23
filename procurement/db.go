package procurement

import "encore.dev/storage/sqldb"

var db = sqldb.NewDatabase("procurement", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
