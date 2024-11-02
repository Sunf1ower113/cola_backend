package composites

import (
	"auth-api/internal/adapters/api"
	apiRecycleBox "auth-api/internal/adapters/api/recycleBox"
	adaptersRecycleBox "auth-api/internal/adapters/db/recycleBox"
	domainRecycleBox "auth-api/internal/domain/recycleBox"
	"database/sql"
)

type RecycleBoxComposite struct {
	Storage domainRecycleBox.RecycleBoxStorage
	Service domainRecycleBox.ServiceRecycleBox
	Handler api.Handler
}

func NewRecycleBoxComposite(db *sql.DB) (*RecycleBoxComposite, error) {
	recycleBoxStorageStorage := adaptersRecycleBox.NewRecycleBoxStorage(db)
	recycleBoxService := domainRecycleBox.NewRecycleBoxService(recycleBoxStorageStorage)
	recycleBoxHandler := apiRecycleBox.NewHandler(recycleBoxService)
	return &RecycleBoxComposite{
		Storage: recycleBoxStorageStorage,
		Service: recycleBoxService,
		Handler: recycleBoxHandler,
	}, nil
}
