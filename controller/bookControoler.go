package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"library_server/model"
	"library_server/response"
	"library_server/service"
)

type BookController struct {
}

// GetBooks
// @Description 查询所有书籍
// @Author John 2023-04-15 15:36:55
// @Param ctx
func (b *BookController) GetBooks(ctx *gin.Context) {
	bookService := service.NewBookService()
	isAdmin := ctx.GetBool("isAdmin")
	books, lErr := bookService.GetBooks(isAdmin)
	// 查询错误
	if lErr != nil {
		fmt.Println(lErr.Err)
		response.Response(ctx, lErr.HttpCode, gin.H{
			"status":     lErr.HttpCode,
			"error_code": lErr.ErrorCode,
			"msg":        lErr.Msg,
		})
		return
	}
	response.Response(ctx, http.StatusOK, gin.H{
		"status":     200,
		"error_code": 1,
		"msg":        "书籍请求成功",
		"data":       books,
	})
}

// GetBooksByName
// @Description 查询书籍
// @Author John 2023-04-18 15:33:55
// @Param ctx
func (b *BookController) GetBooksByName(ctx *gin.Context) {
	bookService := service.NewBookService()
	name := ctx.PostForm("name")
	isAdmin := cast.ToBool(ctx.PostForm("isAdmin"))
	// name为空，跳转到QueryBooks

	fmt.Println("-----2-----------", ctx.GetString("username"))
	fmt.Println("-----2-----------", ctx.GetString("isAdmin"))

	if name == "" {
		books, lErr := bookService.GetBooks(isAdmin)
		if lErr != nil {
			fmt.Println(lErr.Err)
			response.Response(ctx, lErr.HttpCode, gin.H{
				"status":     lErr.HttpCode,
				"error_code": lErr.ErrorCode,
				"msg":        lErr.Msg,
			})
			return
		}
		response.Response(ctx, http.StatusOK, gin.H{
			"status":     200,
			"msg":        "书籍请求成功",
			"error_code": 1,
			"data":       books,
		})
		return
	}

	books, lErr := bookService.GetBookByName(name, isAdmin)
	// 查询出错
	if lErr != nil {
		fmt.Println(lErr.Err)
		response.Response(ctx, lErr.HttpCode, gin.H{
			"status":     lErr.HttpCode,
			"error_code": lErr.ErrorCode,
			"msg":        lErr.Msg,
		})
		return
	}
	response.Success(ctx, gin.H{
		"status":     200,
		"error_code": 1,
		"msg":        "查询成功",
		"data":       books,
	})
}

// UpdateBookInfo
// @Description 管理员更新图书信息
// @Author John 2023-04-27 13:08:35
// @Param ctx
func (b *BookController) UpdateBookInfo(ctx *gin.Context) {
	bookId := ctx.PostForm("bookId")
	value := ctx.PostForm("value")
	status := ctx.PostForm("status")
	difference := ctx.PostForm("difference")

	bookService := service.NewBookService()
	lErr := bookService.UpdateBookInfo(bookId, value, status, difference)
	if lErr != nil {
		fmt.Println(lErr.Err)
		response.Response(ctx, lErr.HttpCode, gin.H{
			"status":     lErr.HttpCode,
			"error_code": lErr.ErrorCode,
			"msg":        lErr.Msg,
		})
		return
	}
	response.Success(ctx, gin.H{
		"status":     200,
		"error_code": 1,
		"msg":        "更新书籍成功",
	})
}

// DeleteBook
// @Description 管理员删除书籍
// @Author John 2023-04-27 20:34:19
// @Param ctx
func (b *BookController) DeleteBook(ctx *gin.Context) {
	// 数据接收
	bookId := ctx.PostForm("bookId")
	bookService := service.NewBookService()
	lErr := bookService.DeleteBook(bookId)

	if lErr != nil {
		fmt.Println(lErr.Err)
		response.Response(ctx, lErr.HttpCode, gin.H{
			"status":     lErr.HttpCode,
			"error_code": lErr.ErrorCode,
			"msg":        lErr.Msg,
		})
		return
	}
	response.Success(ctx, gin.H{
		"status":     200,
		"error_code": 1,
		"msg":        "删除书籍成功",
	})
}

// CreateBook
// @Description 添加图书
// @Author John 2023-05-03 16:27:29
// @Param ctx
func (b *BookController) CreateBook(ctx *gin.Context) {
	bookName := ctx.PostForm("bookName")
	author := ctx.PostForm("author")
	press := ctx.PostForm("press")
	amount := ctx.PostForm("amount")
	position := ctx.PostForm("position")
	Isbn := ctx.PostForm("isbn")

	Amount, err := strconv.Atoi(amount)
	if err != nil {
		fmt.Println("Atoi错误")
		response.Response(ctx, http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "请求错误",
		})
		return
	}
	book := model.Book{
		BookName: bookName,
		Press:    press,
		Isbn:     Isbn,
		Amount:   uint(Amount),
		Author:   author,
		Position: position,
	}

	bookService := service.NewBookService()
	lErr := bookService.CreateBook(book)

	if lErr != nil {
		fmt.Println(lErr.Err)
		response.Response(ctx, lErr.HttpCode, gin.H{
			"status":     lErr.HttpCode,
			"error_code": lErr.ErrorCode,
			"msg":        lErr.Msg,
		})
		return
	}
	response.Success(ctx, gin.H{
		"status":     200,
		"error_code": 1,
		"msg":        "添加图书成功",
	})
}

// NewBookController
// @Description  BookController的构造器
// @Author John 2023-04-16 15:21:28
// @Return BookController
func NewBookController() BookController {
	return BookController{}
}
