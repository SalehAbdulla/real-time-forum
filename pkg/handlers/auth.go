package handlers

import (
	"encoding/json"
	realtimeforum "real-time-forum"
	"net/http"
	"strings"
)

type registerResponse struct {
	Message string `json:"message"`
	UserID  int64  `json:"userId"`
}

type loginResponse struct {
	Message string `json:"message"`
	UserID  int64  `json:"userId"`
}

func (re *Repository) Register(w http.ResponseWriter, r *http.Request) {
	nickName := strings.TrimSpace(strings.ToLower(r.FormValue("nickName")))
	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	firstName := strings.TrimSpace(strings.ToLower(r.FormValue("firstName")))
	lastName := strings.TrimSpace(strings.ToLower(r.FormValue("lastName")))
	password := strings.TrimSpace(r.FormValue("password"))
	confirmPassword := strings.TrimSpace(r.FormValue("confirmPassword"))
	age := strings.TrimSpace(strings.ToLower(r.FormValue("age")))
	gender := strings.TrimSpace(strings.ToLower(r.FormValue("gender")))

	inputs := []string{nickName, email, firstName, lastName, password,
		confirmPassword, age, gender}

	for _, v := range inputs {
		if v == "" {
			re.HandleError(w, r, realtimeforum.ErrBadRequest)
			return
		}
	}

	userID, err := re.AuthService.Register(inputs)
	if err != nil {
		re.App.Logger.Info("registration failed",
			"email", email,
			"error", err.Error(),
		)
		re.HandleError(w, r, err)
		return
	}

	re.App.Logger.Info("user registered successfully",
		"user_id", userID,
		"email", email,
		"nickname", nickName,
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerResponse{
		Message: "user registered successfully",
		UserID:  userID,
	})
}

func (re *Repository) Login(w http.ResponseWriter, r *http.Request) {
	identifier := strings.TrimSpace(strings.ToLower(r.FormValue("identifier")))
	password := strings.TrimSpace(r.FormValue("password"))

	if identifier == "" || password == "" {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	userID, err := re.AuthService.Login(identifier, password)
	if err != nil {
		re.App.Logger.Info("login failed",
			"identifier", identifier,
			"error", err.Error(),
		)
		if err == realtimeforum.ErrInvalidCredentials {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	re.App.Logger.Info("login successful",
		"user_id", userID,
		"identifier", identifier,
	)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{
		Message: "login successful",
		UserID:  userID,
	})
}
