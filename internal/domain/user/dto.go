package user

type CreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserDTO struct {
	ID          int64  `json:"user_id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birth_date"`
	Points      int64  `json:"points"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
}

type AuthDTO struct {
	ID             int64  `json:"user_id"`
	HashedPassword string `json:"hashed_password"`
	Role           string `json:"role"`
}
