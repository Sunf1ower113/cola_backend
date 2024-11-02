package user

type User struct {
	ID             int64  `json:"user_id"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	HashedPassword string `json:"password"`
	PhoneNumber    string `json:"phone_number"`
	BirthDate      string `json:"birth_date"`
	Points         int64  `json:"points"`
	Role           string `json:"role"`
}
