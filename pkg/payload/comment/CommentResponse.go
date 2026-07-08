package comment

type CommentResponse struct {
	CommentDTO    []CommentDTO `json:"posts"`
	PageNumber    int          `json:"pageNumber"`
	PageSize      int          `json:"pageSize"`
	TotalElements int          `json:"totalElements"`
	TotalPages    int          `json:"totalPages"`
	LastPage      bool         `json:"lastPage"`
}
