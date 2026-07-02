package handlers

import (
	"net/http"
	realtimeforum "real-time-forum"
)

func (re *Repository) parseForm(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return false 
	}
	return true
}

func (re *Repository) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	// Does Nothing
	
}
