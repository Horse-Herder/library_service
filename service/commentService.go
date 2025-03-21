package service

import (
	"errors"
	"gorm.io/gorm"
	"library_server/common"
	"library_server/model"
	"library_server/repository"
	"library_server/utils"
	"library_server/vo"
	"net/http"
)

type CommentService struct {
	DB *gorm.DB
}

// GetComments
// @Description 查询所有评论
// @Author John 2023-04-20 21:10:09
// @Return []vo.CommentVO
// @Return *common.LError
func (c *CommentService) GetComments() (comments []vo.CommentVo, lErr *common.LError) {
	//c.DB.Preload("Reader").Preload("Book").Find(&comments)
	commentRepository := repository.NewCommentRepository()
	comments, err := commentRepository.GetCommentVos()
	// 获取评论CommentVo
	if err != nil {
		return comments, &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "查询所有评论失败",
			Err:      errors.New("获取评论CommentVo失败"),
		}
	}
	return comments, nil
}

// GetCommentAmount
// @Description 返回评论数量
// @Author John 2023-04-25 20:04:45
// @Return amount
// @Return lErr
func (c *CommentService) GetCommentAmount() (count int64, lErr *common.LError) {
	commentRepository := repository.NewCommentRepository()
	count, err := commentRepository.GetCommentCount()

	//  查询评论数量
	if err != nil {
		return count, &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "查询评论数量错误",
			Err:      errors.New("查询评论数量失败"),
		}
	}
	return count, nil
}

// CreateComment
// @Description 添加评论
// @Author John 2023-04-25 22:06:46
// @Param readerId
// @Param bookId
// @Param content
// @Return lErr
func (c *CommentService) CreateComment(readerId string, bookId string, content string) (lErr *common.LError) {
	// 数据验证
	if readerId == "" || bookId == "" || content == "" {
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "添加评论失败",
			Err:      errors.New("数据验证失败"),
		}
	}

	commentRepository := repository.NewCommentRepository()
	comment := model.Comment{
		ReaderId: readerId,
		BookId:   bookId,
		Date:     model.Time(utils.NowTime()),
		Content:  content,
		Praise:   0,
		Status:   1,
	}
	// 开启事务
	tx := c.DB.Begin()
	// 添加评论
	err := commentRepository.CreateComment(tx, comment)
	if err != nil {
		tx.Rollback()
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "添加评论失败",
			Err:      errors.New("添加评论失败"),
		}
	}
	tx.Commit()
	return nil
}

// UpdatePraise
// @Description 更新点赞记录
// @Author John 2023-04-28 16:19:30
// @Param readerId
// @Param bookId
// @Param date
// @Return lErr
func (c *CommentService) UpdatePraise(readerId string, bookId string, date string) (lErr *common.LError) {
	// 数据验证
	if readerId == "" || bookId == "" || date == "" {
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "更新点赞记录失败",
			Err:      errors.New("数据验证失败"),
		}
	}
	t, _ := utils.ParseTime(date)
	time := model.Time(t)

	commentRepository := repository.NewCommentRepository()
	// 获取commentId
	commentId, err := commentRepository.GetCommentId(readerId, bookId, time)
	if err != nil {
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "更新点赞记录失败",
			Err:      errors.New("获取commentId失败"),
		}
	}

	// 开启事务
	tx := c.DB.Begin()
	// 添加评论
	err = commentRepository.UpdatePraiseByCommentId(tx, commentId)
	if err != nil {
		tx.Rollback()
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "添加评论失败",
			Err:      errors.New("添加评论失败"),
		}
	}
	tx.Commit()
	return nil
}

func NewCommentService() CommentService {
	return CommentService{
		DB: common.GetDB(),
	}
}
