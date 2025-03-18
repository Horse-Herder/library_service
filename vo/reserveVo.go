package vo

import "library_server/model"

type ReserveVo struct {
	ReaderId   string     `json:"readerId"`
	BookId     string     `json:"bookId"`
	Status     string     `json:"status"`
	Author     string     `json:"author"`
	BookName   string     `json:"bookName"`
	Date       model.Time `json:"date"`
	ReaderName string     `json:"readerName"`
}
