package error

import "errors"

const (
	LoginUserErrorMsg           = "invalid email or password"
	CreateUserBadInputErrorMsg  = "invalid registration data"
	UpdateUserBadInputErrorMsg  = "invalid update data"
	NothingToUpdateUserErrorMsg = "nothing to update"
	BusyUpdateEmailErrorMsg     = "email is busy"
	UserNotFoundErrorMsg        = "not found"
	BoxFullErrorMsg             = "recycle box is full"
)

var (
	NotFoundError           = errors.New(UserNotFoundErrorMsg)
	NothingToUpdateError    = errors.New(NothingToUpdateUserErrorMsg)
	LoginError              = errors.New(LoginUserErrorMsg)
	BusyUpdateEmailError    = errors.New(BusyUpdateEmailErrorMsg)
	CreateUserBadInputError = errors.New(CreateUserBadInputErrorMsg)
	UpdateUserBadInputError = errors.New(UpdateUserBadInputErrorMsg)
	BoxFullError            = errors.New(BoxFullErrorMsg) // Новая ошибка для полной корзины
)
