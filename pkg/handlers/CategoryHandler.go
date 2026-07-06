package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-forum/pkg/models"
)

func (re *Repository) GetCategories(w http.ResponseWriter, r *http.Request) {
	response, err := re.CategoryService.GetCategories()

	if err != nil {
		re.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.DataResponse[[]models.Category]{
		Data: response,
	})
}
