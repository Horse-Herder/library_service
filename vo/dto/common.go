package dto

type PageQuery struct {
	Page  int `json:"page" binding:"min=1"`
	Limit int `json:"limit" binding:"max=1000"`
}
