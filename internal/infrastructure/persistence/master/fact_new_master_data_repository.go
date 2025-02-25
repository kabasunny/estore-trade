// internal/infrastructure/persistence/master/fact_new_master_data_repository.go
package master

import (
	"database/sql"
	"estore-trade/internal/domain"
)

func NewMasterDataRepository(db *sql.DB) domain.MasterDataRepository {
	return &masterDataRepository{db: db}
}
