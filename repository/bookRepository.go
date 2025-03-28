package repository

import (
	"log"

	"gorm.io/gorm"

	"library_server/common"
	"library_server/model"
	"library_server/vo/dto"
)

type BookRepository struct {
	DB *gorm.DB
}

// GetBooks
// @Description 查询所有书籍
func (b *BookRepository) GetBooks(isAdmin bool, req *dto.BookListRequest) (books []model.Book, total int64, err error) {

	if req == nil {
		req.PageQuery = &dto.PageQuery{
			Page:  1,
			Limit: 10,
		}
	}
	name := req.Name
	// 计算偏移量
	offset := (req.PageQuery.Page - 1) * req.PageQuery.Limit

	// 构建基础查询
	query := b.DB.Model(&books)
	query = query.Where("is_deleted = ?", 0)

	// 根据 isAdmin 条件添加查询条件
	if !isAdmin {
		query = query.Where("status = ? ", 1)
	}

	if req.Name != "" {
		query = query.Where("book_name LIKE ? OR author LIKE ?", "%"+name+"%", "%"+name+"%")
	}
	
	// 查询总数
	if err = query.Count(&total).Error; err != nil {
		log.Printf("Failed to count books: %v", err)
		return nil, 0, err
	}

	// 应用分页并查询数据
	result := query.Offset(offset).Limit(req.Limit).Find(&books)
	if result.Error != nil {
		log.Printf("Failed to query books: %v", result.Error)
		return nil, 0, result.Error
	}

	return books, total, nil
}

// GetBooksByName
// @Description 根据书名查询书籍
func (b *BookRepository) GetBooksByName(bookName string, isAdmin bool) (books []model.Book, err error) {
	if isAdmin {
		if err := b.DB.Where("book_name like ?", "%"+bookName+"%", "is_deleted = ?", 0).Find(&books).Error; err != nil {
			return books, err
		}
	} else {
		if err := b.DB.Where("book_name like ?", "%"+bookName+"%", "status = ?", 1, "is_deleted = ?", 0).Find(&books).Error; err != nil {
			return books, err
		}
	}
	return books, nil
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
	// if err := tx.Where("book_id = ?", bookId).Delete(&model.Book{}).Error; err != nil {
	// 	return err
	// }
	if err := tx.Model(&model.Book{}).Where("book_id = ?", bookId).UpdateColumn("is_deleted", 1).Error; err != nil {
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
