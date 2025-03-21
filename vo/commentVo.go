package vo

import "library_server/model"

type CommentVo struct {
	Email      string     `json:"email"`
	Status     uint       `json:"status"`
	CommentId  string     `json:"commentId"`
	ReaderId   string     `json:"readerId"`
	BookId     string     `json:"bookId"`
	ReaderName string     `json:"readerName"`
	BookName   string     `json:"bookName"`
	Date       model.Time `json:"date"`
	Content    string     `json:"content"`
	Praise     uint       `json:"praise"`
}
