package service

import (
	"net/mail"
	realtimeforum "real-time-forum"
	"strconv"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	db "real-time-forum/pkg/repos"
)

type AuthService interface {
	Register(inputs []string) (int64, error)
	Login(identifier, password string) (int64, error)
}

type authServiceImpl struct {
	db *db.DB
}

func NewAuthService(database *db.DB) AuthService {
	return authServiceImpl{
		db: database,
	}
}

// passwordStrength checks that the password contains at least one letter,
// one number, and one symbol.
func passwordStrength(password string) bool {
	var hasLetter, hasNumber, hasSymbol bool
	for _, ch := range password {
		switch {
		case unicode.IsLetter(ch):
			hasLetter = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsSymbol(ch) || unicode.IsPunct(ch):
			hasSymbol = true
		}
	}
	return hasLetter && hasNumber && hasSymbol
}

func (s authServiceImpl) Register(inputs []string) (int64, error) {
	if len(inputs) < 8 {
		return 0, realtimeforum.ErrBadRequest
	}

	nickName := inputs[0]
	email := inputs[1]
	firstName := inputs[2]
	lastName := inputs[3]
	password := inputs[4]
	confirmPassword := inputs[5]
	ageStr := inputs[6]
	gender := inputs[7]

	if len(nickName) < 2 || len(nickName) > 33 {
		return 0, realtimeforum.ErrNickNameLength
	}

	// Validate email format using net/mail
	_, err := mail.ParseAddress(email)
	if err != nil {
		return 0, realtimeforum.ErrInvalidEmail
	}

	// Check if email already exists
	if err := s.db.DoesEmailExists(email); err != nil {
		return 0, err
	}

	// Check if nickname already exists
	if err := s.db.DoesNicknameExists(nickName); err != nil {
		return 0, err
	}

	// Validate age and calculate yearOfBirth
	age, err := strconv.Atoi(ageStr)
	if err != nil || age <= 0 || age > 150 {
		return 0, realtimeforum.ErrInvalidAge
	}
	yearOfBirth := time.Now().Year() - age

	// Confirm passwords match
	if password != confirmPassword {
		return 0, realtimeforum.ErrPasswordsDontMatch
	}

	if len(password) < 12 || len(password) > 64 {
		return 0, realtimeforum.ErrPasswordLength
	}


	// Validate password strength
	if !passwordStrength(password) {
		return 0, realtimeforum.ErrInvalidPassForm
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, realtimeforum.ErrInternal
	}

	if gender != "male" && gender != "female" {
		return 0, realtimeforum.ErrGender
	}

	// Insert user into the database via repository
	return s.db.InsertUser(nickName, firstName, lastName, email, string(hashedPassword), yearOfBirth, gender)
}

func (s authServiceImpl) Login(identifier, password string) (int64, error) {
	userID, hashedPassword, err := s.db.GetUserCredentials(identifier)
	if err != nil {
		return 0, err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return 0, realtimeforum.ErrInvalidCredentials
	}

	return userID, nil
}
