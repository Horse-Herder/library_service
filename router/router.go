package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"library_server/controller"
	"library_server/middleware"
)

type RATPConsoleServer struct {
	engine *gin.Engine
	server *http.Server

	forever chan os.Signal
}

func CollectRoute(r *gin.Engine) *gin.Engine {
	// 配置CORS跨域路由
	CORSMiddleware := middleware.CORSMiddleware()
	base := r.Group("/")

	r.Use(CORSMiddleware)

	userController := controller.NewUserController()
	r.POST("/login", userController.Login)
	r.POST("/register", userController.Register)
	r.POST("/logout", userController.Logout)

	base.Use(JWT())

	// book
	bookController := controller.NewBookController()
	base.POST("/books", bookController.GetBooks)
	//base.POST("/searchbook", bookController.GetBooksByName)
	base.POST("/changebookinfo", bookController.UpdateBookInfo)
	base.POST("/delbook", bookController.DeleteBook)
	base.POST("/adminaddbooks", bookController.CreateBook)

	// comment
	commentController := controller.NewCommentController()
	base.POST("/comments", commentController.GetComments)
	base.POST("/amount", commentController.GetCommentCount)
	base.POST("/addcomment", commentController.CreateComment)
	base.POST("/addpraise", commentController.UpdatePraise)

	// reader
	readerController := controller.NewReaderController()
	base.POST("/initreader", readerController.GetReaderInfo)
	base.POST("/amountmax", readerController.GetMaxCountReader)
	base.POST("/delperson", readerController.DeleteReader)
	base.POST("/initreaderlist", readerController.GetReaders)

	// borrow
	borrowController := controller.NewBorrowController()
	base.POST("/addborrow", borrowController.CreateBorrowRecord)
	base.POST("/borrows", borrowController.GetReaderBorrowRecords)
	base.POST("/returnbook", borrowController.ReturnBook)
	base.POST("/continueborrow", borrowController.RenewBook)
	base.POST("/borrowslist", borrowController.GetAllBorrowRecords)
	base.POST("/searchborrow", borrowController.GetBorrowRecordByInfo)
	base.POST("/deleteborrow", borrowController.DeleteBorrow)
	base.POST("/alertperson", borrowController.SendReminder)

	// reserve
	reserveController := controller.NewReserveController()
	base.POST("/addreserve", reserveController.CreateReserveRecord)
	base.POST("/reserve", reserveController.GetReserveRecords)
	base.POST("/cancelreserve", reserveController.DeleteReserveRecord)
	base.POST("/reservelist", reserveController.GetAllReserveRecords)

	// report
	reportController := controller.NewReportController()
	base.POST("/initstureport", reportController.GetReportRecords)
	base.POST("/initreportlist", reportController.GetAllReportRecords)
	base.POST("/reportcomment", reportController.CreateReport)
	base.POST("/auditcomment", reportController.ManageReport)

	return r
}
