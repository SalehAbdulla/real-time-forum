package user

import (
	"net/http"
	realtimeforum "real-time-forum"
	"net/mail"
	"strconv"
	"strings"
	"unicode"
)

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

func (d *RegisterRequestDTO) ParseAndValidate(r *http.Request) error {
	d.Nickname = strings.TrimSpace(strings.ToLower(r.FormValue("nickName")))
	d.Email = strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	d.FirstName = strings.TrimSpace(strings.ToLower(r.FormValue("firstName")))
	d.LastName = strings.TrimSpace(strings.ToLower(r.FormValue("lastName")))
	d.Password = strings.TrimSpace(r.FormValue("password"))
	d.ConfirmPassword = strings.TrimSpace(r.FormValue("confirmPassword"))
	d.Age = strings.TrimSpace(strings.ToLower(r.FormValue("age")))
	d.Gender = strings.TrimSpace(strings.ToLower(r.FormValue("gender")))

	if d.Nickname == "" || d.Email == "" || d.FirstName == "" || d.LastName == "" ||
		d.Password == "" || d.ConfirmPassword == "" || d.Age == "" || d.Gender == "" {
		return realtimeforum.ErrBadRequest
	}

	if len(d.Nickname) < 2 || len(d.Nickname) > 33 {
		return realtimeforum.ErrNickNameLength
	}

	if _, err := mail.ParseAddress(d.Email); err != nil {
		return realtimeforum.ErrInvalidEmail
	}

	if d.Password != d.ConfirmPassword {
		return realtimeforum.ErrPasswordsDontMatch
	}

	if len(d.Password) < 12 || len(d.Password) > 64 {
		return realtimeforum.ErrPasswordLength
	}

	if !passwordStrength(d.Password) {
		return realtimeforum.ErrInvalidPassForm
	}

	age, err := strconv.Atoi(d.Age)
	if err != nil || age <= 0 || age > 100 {
		return realtimeforum.ErrInvalidAge
	}

	if d.Gender != "male" && d.Gender != "female" {
		return realtimeforum.ErrGender
	}

	return nil
}