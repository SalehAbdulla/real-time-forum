package service

import (
	"log/slog"
	"net/mail"
	realtimeforum "real-time-forum"
	"strconv"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"real-time-forum/pkg/models"
	db "real-time-forum/pkg/repositories"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(inputs models.RegisterRequest) (string, string, error)
	Login(identifier, password string) (string, string, error)
	Logout(token string) error
	GetMe(userID string) (models.UserProfile, error)
}

type authServiceImpl struct {
	db             db.AuthRepository
	sessionManager *SessionManager
}

func NewAuthService(database db.AuthRepository) AuthService {
	return authServiceImpl{
		db:             database,
		sessionManager: NewSessionManager(),
	}
}

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

func (s authServiceImpl) Register(registerRequest models.RegisterRequest) (string, string, error) {

	nickName := registerRequest.Nickname
	email := registerRequest.Email
	firstName := registerRequest.FirstName
	lastName := registerRequest.LastName
	password := registerRequest.Password
	confirmPassword := registerRequest.ConfirmPassword
	ageStr := registerRequest.Age
	gender := registerRequest.Gender

	if len(nickName) < 2 || len(nickName) > 33 {
		slog.Debug("invalid nickname length", "nickname", nickName, "length", len(nickName))
		return "", "", realtimeforum.ErrNickNameLength
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		slog.Debug("invalid email format", "email", email)
		return "", "", realtimeforum.ErrInvalidEmail
	}

	if err := s.db.DoesEmailExists(email); err != nil {
		return "", "", err
	}

	if err := s.db.DoesNicknameExists(nickName); err != nil {
		return "", "", err
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil || age <= 0 || age > 150 {
		slog.Debug("invalid age", "age", ageStr)
		return "", "", realtimeforum.ErrInvalidAge
	}
	yearOfBirth := time.Now().Year() - age

	if password != confirmPassword {
		slog.Debug("passwords do not match")
		return "", "", realtimeforum.ErrPasswordsDontMatch
	}

	if len(password) < 12 || len(password) > 64 {
		slog.Debug("invalid password length", "length", len(password))
		return "", "", realtimeforum.ErrPasswordLength
	}

	if !passwordStrength(password) {
		slog.Debug("weak password")
		return "", "", realtimeforum.ErrInvalidPassForm
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return "", "", realtimeforum.ErrInternal
	}

	if gender != "male" && gender != "female" {
		slog.Debug("invalid gender", "gender", gender)
		return "", "", realtimeforum.ErrGender
	}

	userID := uuid.NewString()
	if err := s.db.InsertUser(userID, nickName, firstName, lastName, email, string(hashedPassword), yearOfBirth, gender); err != nil {
		return "", "", err
	}

	token := uuid.NewString()
	s.sessionManager.CreateSession(userID, token)

	return userID, token, nil
}

func (s authServiceImpl) Login(identifier, password string) (string, string, error) {
	userID, hashedPassword, err := s.db.GetUserCredentials(identifier)
	if err != nil {
		slog.Info("login credential lookup failed", "identifier", identifier, "error", err)
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Info("invalid password attempt", "user_id", userID)
		return "", "", realtimeforum.ErrInvalidCredentials
	}

	token := uuid.NewString()
	s.sessionManager.CreateSession(userID, token)

	return userID, token, nil
}

func (s authServiceImpl) Logout(token string) error {
	s.sessionManager.DeleteSession(token)
	return nil
}

func (s authServiceImpl) GetMe(userID string) (models.UserProfile, error) {
	return s.db.GetUserProfile(userID)
}