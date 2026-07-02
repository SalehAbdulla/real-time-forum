package handlers

import (
	"net/http"
	"strings"
)

func (re *Repository) Register(w http.ResponseWriter, r *http.Request) {

	nickName := strings.TrimSpace(strings.ToLower(r.FormValue("nickName")))
	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	firstName := strings.TrimSpace(strings.ToLower(r.FormValue("firstName")))
	lastName := strings.TrimSpace(strings.ToLower(r.FormValue("lastName")))
	password := strings.TrimSpace(strings.ToLower(r.FormValue("password")))
	confirmPassword := strings.TrimSpace(strings.ToLower(r.FormValue("confirmPassword")))
	age := strings.TrimSpace(strings.ToLower(r.FormValue("age")))
	gender := strings.TrimSpace(strings.ToLower(r.FormValue("gender")))

	inputDictionary := []string{nickName, email, firstName, lastName, password,
		confirmPassword, age, gender}

	for _, v := range inputDictionary {
		if v == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	

}
