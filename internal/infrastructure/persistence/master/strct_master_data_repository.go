// internal/infrastructure/persistence/master/strct_master_data_repository.go
package master

import "database/sql"

type masterDataRepository struct {
	db *sql.DB
}
