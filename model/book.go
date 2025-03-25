package model

type Book struct {
	BookId   string `json:"bookId" gorm:"type:varchar(50);primary_key;"`
	BookName string `json:"bookName" gorm:"type:varchar(20);"`
	Press    string `json:"press" gorm:"type:varchar(20);"`
	Author   string `json:"author" gorm:"varchar(10);not null"`
	Isbn     string `json:"isbn" gorm:"varchar(10);not null"`
	// 当前数量
	Amount uint `json:"amount"`
	// 位置
	Position string `json:"position" gorm:"type:varchar(30);"`
	// 总数量
	TotalAmount uint `json:"totalAmount"`
	// 借阅次数
	BorrowedTimes uint `json:"borrowedTimes"`
	// 状态 1:可借阅 0：不可借阅
	Status int `json:"status" gorm:"type:int(10);"`
	// 是否删除 0:未删除 1:已删除
	IsDeleted int `json:"is_deleted" gorm:"type:int(10);"`
}
