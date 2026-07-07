package posts

type PostResponse struct {
	Posts         []PostDTO
	pageNumber    int
	pageSize      int
	totalElements int
	totalPages    int
	lastPage      bool
}
