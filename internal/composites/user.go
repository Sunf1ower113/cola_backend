package composites

import (
	"auth-api/internal/adapters/api"
	apiUser "auth-api/internal/adapters/api/user"
	adaptersUser "auth-api/internal/adapters/db/user"
	domainUser "auth-api/internal/domain/user"
	"database/sql"
)

type UserComposite struct {
	Storage domainUser.StorageUser
	Service domainUser.ServiceUser
	Handler api.Handler
}

func NewUserComposite(db *sql.DB) (*UserComposite, error) {
	userStorage := adaptersUser.NewUserStorage(db)
	userService := domainUser.NewUserService(userStorage)
	userHandler := apiUser.NewHandler(userService)
	return &UserComposite{
		Storage: userStorage,
		Service: userService,
		Handler: userHandler,
	}, nil
}
