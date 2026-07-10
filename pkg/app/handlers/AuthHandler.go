package handlers

import (
	"encoding/json"
	"net/http"
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/middleware"
	"real-time-forum/pkg/payload"
	"real-time-forum/pkg/payload/user"
	"strings"
)

type registerResponse struct {
	UserID   string `json:"userId"`
	Nickname string `json:"nickname"`
}

type loginResponse struct {
	UserID   string `json:"userId"`
	Nickname string `json:"nickname"`
}

type logoutResponse struct {
	Message string `json:"message"`
}

func (re *HandlerContext) Register(w http.ResponseWriter, r *http.Request) {
	registerRequest := user.RegisterRequestDTO{}

	registerRequest.Nickname = strings.TrimSpace(strings.ToLower(r.FormValue("nickName")))
	registerRequest.Email = strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	registerRequest.FirstName = strings.TrimSpace(strings.ToLower(r.FormValue("firstName")))
	registerRequest.LastName = strings.TrimSpace(strings.ToLower(r.FormValue("lastName")))
	registerRequest.Password = strings.TrimSpace(r.FormValue("password"))
	registerRequest.ConfirmPassword = strings.TrimSpace(r.FormValue("confirmPassword"))
	registerRequest.Age = strings.TrimSpace(strings.ToLower(r.FormValue("age")))
	registerRequest.Gender = strings.TrimSpace(strings.ToLower(r.FormValue("gender")))

	if registerRequest.Nickname == "" || registerRequest.Email == "" ||
		registerRequest.FirstName == "" || registerRequest.LastName == "" ||
		registerRequest.Password == "" || registerRequest.ConfirmPassword == "" ||
		registerRequest.Age == "" || registerRequest.Gender == "" {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return
	}

	userID, token, err := re.AuthService.Register(registerRequest)
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
		"email", registerRequest.Email,
		"nickname", registerRequest.Nickname,
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payload.SuccessResponse[registerResponse]{
		Success: true,
		Data: registerResponse{
			UserID:   userID,
			Nickname: registerRequest.Nickname,
		},
		Message: "user registered successfully",
	})

}

func (re *HandlerContext) Login(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(payload.SuccessResponse[loginResponse]{
		Success: true,
		Data: loginResponse{
			UserID: userID,
		},
		Message: "login successful",
	})
}

func (re *HandlerContext) Logout(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(payload.SuccessResponse[logoutResponse]{
		Success: true,
		Data: logoutResponse{
			Message: "Logged out",
		},
		Message: "Logged out successfully",
	})
}

func (re *HandlerContext) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
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
	json.NewEncoder(w).Encode(payload.SuccessResponse[user.UserDTO]{
		Success: true,
		Data:    profile,
		Message: "User profile retrieved successfully",
	})
}