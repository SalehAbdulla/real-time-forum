package user

type RegisterRequestDTO struct {
	Nickname        string
	Email           string
	FirstName       string
	LastName        string
	Password        string
	ConfirmPassword string
	Age             string
	Gender          string
}