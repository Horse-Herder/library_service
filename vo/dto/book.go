package dto

type BookListRequest struct {
	Name string `json:"name"`
	*PageQuery
}
