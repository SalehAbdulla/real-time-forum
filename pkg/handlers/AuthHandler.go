package handlers

import (
	"encoding/json"
	"net/http"
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/models"
	"strings"
)

type registerResponse struct {
	Message string `json:"message"`
	UserID  string `json:"userId"`
}

type loginResponse struct {
	Message string `json:"message"`
	UserID  string `json:"userId"`
}

type dataResponse[T any] struct {
	Data T `json:"data"`
}

type logoutResponse struct {
	Message string `json:"message"`
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

	userID, token, err := re.AuthService.Register(inputs)
    if err != nil {
        re.HandleError(w, r, err)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        Secure:   re.App.InProduction,
        SameSite: http.SameSiteStrictMode,
    })

    re.App.Logger.Info("user registered and logged in successfully",
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

	userID, token, err := re.AuthService.Login(identifier, password)
	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   re.App.InProduction,
		SameSite: http.SameSiteStrictMode,
	})

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

func (re *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("session_token")
	if err != nil || tokenCookie.Value == "" {
		re.HandleError(w, r, realtimeforum.ErrUnauthorized)
		return
	}

	if err := re.AuthService.Logout(tokenCookie.Value); err != nil {
		re.HandleError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   re.App.InProduction,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataResponse[logoutResponse]{
		Data: logoutResponse{Message: "Logged out"},
	})
}

func (re *Repository) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		re.HandleError(w, r, realtimeforum.ErrUnauthorized)
		return
	}

	profile, err := re.AuthService.GetMe(userID)
	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataResponse[models.UserProfile]{
		Data: profile,
	})
}
