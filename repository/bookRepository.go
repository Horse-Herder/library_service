package repository

import (
	"gorm.io/gorm"

	"library_server/common"
	"library_server/model"
)

type BookRepository struct {
	DB *gorm.DB
}

// GetBooks
// @Description 查询所有书籍
func (b *BookRepository) GetBooks(isAdmin bool, page int, pageSize int) (books []model.Book, total int64, err error) {
	// Build the base query
	query := b.DB.Model(&model.Book{})

	if !isAdmin {
		query = query.Where("status = ?", 1)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil

}

// GetBooksByName
// @Description 根据书名查询书籍
func (b *BookRepository) GetBooksByName(searchName string, isAdmin bool, page int, pageSize int) (books []model.Book, total int64, err error) {
	query := b.DB.Model(&model.Book{})

	// 模糊搜索书名或作者
	if searchName != "" {
		query = query.Where("book_name LIKE ? OR author LIKE ?", "%"+searchName+"%", "%"+searchName+"%")
	}

	// 如果不是管理员，则只显示 status = 1 的书籍
	if !isAdmin {
		query = query.Where("status = ?", 1)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// UpdateBookAmount
// @Description 更新书籍总数
func (b *BookRepository) UpdateBookAmount(tx *gorm.DB, bookId string, count int) error {
	return tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("amount", gorm.Expr("amount + ?", count)).Error
}

// UpdateBookBorrowedTimes
// @Description 更新书籍借阅次数
func (b *BookRepository) UpdateBookBorrowedTimes(tx *gorm.DB, bookId string, count int) error {
	return tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("borrowed_times", gorm.Expr("borrowed_times + ?", count)).Error
}

// UpdateBookNameByBookId
// @Description 更新书名
func (b *BookRepository) UpdateBookNameByBookId(tx *gorm.DB, bookId string, bookName string) error {
	if err := tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("book_name", bookName).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAuthorByBookId
// @Description 更新作者
func (b *BookRepository) UpdateAuthorByBookId(tx *gorm.DB, bookId string, author string) interface{} {
	if err := tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("author", author).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePositionByBookId
// @Description  更新书籍位置
func (b *BookRepository) UpdatePositionByBookId(tx *gorm.DB, bookId string, position string) error {
	if err := tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("position", position).Error; err != nil {
		return err
	}
	return nil
}

// GetBookByPosition
// @Description 返回指定位置的图书
func (b *BookRepository) GetBookByPosition(position string) (book model.Book, err error) {
	if err = b.DB.Model(&model.Book{}).Where("position = ?", position).First(&book).Error; err != nil {
		return book, err
	}
	return book, nil
}

// UpdateTotalAmountByBookId
// @Description 更新总数量
func (b *BookRepository) UpdateTotalAmountByBookId(tx *gorm.DB, bookId string, count int) error {
	if err := tx.
		Model(&model.Book{}).
		Where("book_id = ?", bookId).
		UpdateColumn("amount", gorm.Expr("amount + ?", count)).
		Error; err != nil {
		return err
	}
	return nil
}

// UpdateAmountByBookId
// @Description 更新当前数量
func (b *BookRepository) UpdateAmountByBookId(tx *gorm.DB, bookId string, count int) error {
	if err := tx.
		Model(&model.Book{}).
		Where("book_id = ?", bookId).
		UpdateColumn("total_amount", gorm.Expr("total_amount + ?", count)).
		Error; err != nil {
		return err
	}
	return nil
}

// UpdatePressByBookId
// @Description 出版社
func (b *BookRepository) UpdatePressByBookId(tx *gorm.DB, bookId string, press string) interface{} {
	if err := tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("press", press).Error; err != nil {
		return err
	}
	return nil
}

// UpdateISBNByBookId
// @Description ISBN
func (b *BookRepository) UpdateISBNByBookId(tx *gorm.DB, bookId string, Isbn string) interface{} {
	if err := tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("isbn", Isbn).Error; err != nil {
		return err
	}
	return nil
}

// UpdateStatusByBookId
// @Description Status
func (b *BookRepository) UpdateStatusByBookId(tx *gorm.DB, bookId string, status int64) interface{} {
	if err := tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBookByBookId
// @Description 根据书籍id删除书籍
func (b *BookRepository) DeleteBookByBookId(tx *gorm.DB, bookId string) error {
	if err := tx.Where("book_id = ?", bookId).Delete(&model.Book{}).Error; err != nil {
		return err
	}
	return nil
}

// GetAmountByBookId
// @Description 返回当前书籍当前库存
func (b *BookRepository) GetAmountByBookId(bookId string) (amount int, err error) {
	if err = b.DB.Model(&model.Book{}).Select(`amount`).Where("book_id = ?", bookId).Scan(&amount).Error; err != nil {
		return amount, err
	}
	return amount, nil
}

// GetTotalAmountByBookId
// @Description  返回当前书籍总库存
func (b *BookRepository) GetTotalAmountByBookId(bookId string) (totalAmount int, err error) {
	if err = b.DB.Model(&model.Book{}).Select(`total_amount`).Where("book_id = ?", bookId).Scan(&totalAmount).Error; err != nil {
		return totalAmount, err
	}
	return totalAmount, nil
}

// GetBookIdByBookName
// @Description 根据书籍名称获取书籍id
func (b *BookRepository) GetBookIdByBookName(bookName string) (bookId string, err error) {
	if err = b.DB.Model(&model.Book{}).Select(`book_id`).Where("book_name = ?", bookName).Scan(&bookId).Error; err != nil {
		return bookId, err
	}
	return bookId, nil
}

// CreateBook
// @Description 新增书籍
func (b *BookRepository) CreateBook(tx *gorm.DB, book model.Book) error {
	if err := tx.Create(&book).Error; err != nil {
		return err
	}

	return nil
}

func NewBookRepository() BookRepository {
	return BookRepository{
		DB: common.GetDB(),
	}
}
