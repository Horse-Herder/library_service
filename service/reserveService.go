package service

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"library_server/common"
	"library_server/model"
	"library_server/repository"
	"library_server/vo"
)

type ReserveService struct {
	DB *gorm.DB
}

// CreateReserveRecord
// @Description 新增预约记录
// @Author John 2023-04-20 22:05:10
// @Param reserve
// @Return lErr
func (r *ReserveService) CreateReserveRecord(addReserve model.Reserve) (lErr *common.LError) {
	reserveRepository := repository.NewReserveRepository()
	if addReserve.ReaderId == "" || addReserve.BookId == "" {
		fmt.Println("预约失败")
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "预约失败",
			Err:      errors.New("预约失败"),
		}
	}
	// 获取id
	id, err := reserveRepository.GetReserveId(addReserve.ReaderId, addReserve.BookId, addReserve.Date)
	if err != nil {
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "新增预约记录失败",
			Err:      errors.New("获取id失败"),
		}
	}

	// 验证数据库是否已经存在该预约
	//reserve, _ := reserveRepository.GetReserveByReaderIDAndBookID(addReserve.ReaderId, addReserve.BookId)
	if id != "" {
		fmt.Println("预约记录已存在")
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "预约记录已存在",
			Err:      errors.New("预约记录已存在"),
		}
	}

	tx := r.DB.Begin()
	if err := reserveRepository.CreateReserveRecord(tx, addReserve); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "新增预约记录失败",
			Err:      errors.New("新增预约记录失败"),
		}
	}
	tx.Commit()
	return nil
}

// GetReserves
// @Description 根据readerId获取预约信息
// @Author John 2023-04-20 22:52:29
// @Param readerId
// @Return reserveVOs
// @Return lErr
func (r *ReserveService) GetReserves(readerId string) (reserveVos []vo.ReserveVo, lErr *common.LError) {
	var reserveRepository = repository.NewReserveRepository()
	if readerId == "" {
		return reserveVos, &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "查询预约记录失败",
			Err:      errors.New("readerId为空"),
		}
	}

	reserveVos, err := reserveRepository.GetReserveVosByReaderId(readerId)
	if err != nil {
		return reserveVos, &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "读者请求预约记录失败",
			Err:      errors.New("读者请求预约记录失败"),
		}
	}

	//查询数据为空
	//if len(reserveVos) == 0 {
	//	fmt.Println("读者请求预约记录为空")
	//	return reserveVos, &common.LError{
	//		HttpCode: http.StatusOK,
	//		Msg:      "读者请求预约记录为空",
	//		Err:      errors.New("读者请求预约记录为空"),
	//	}
	//}
	return reserveVos, nil
}

// DeleteReserveRecord
// @Description 删除预约记录
// @Author John 2023-04-20 22:59:06
// @Param delReserve
// @Return lErr
func (r *ReserveService) DeleteReserveRecord(delReserve model.Reserve) (lErr *common.LError) {
	reserveRepository := repository.NewReserveRepository()
	// 获取id
	id, err := reserveRepository.GetReserveId(delReserve.ReaderId, delReserve.BookId, delReserve.Date)
	if err != nil {
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "新增预约记录失败",
			Err:      errors.New("获取id失败"),
		}
	}
	tx := r.DB.Begin()
	if err := reserveRepository.DeleteReserveRecordById(tx, id); err != nil {
		tx.Rollback()
		return &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "取消预约失败",
			Err:      errors.New("取消预约失败"),
		}
	}
	tx.Commit()
	return nil
}

// GetAllReserveRecords
// @Description 管理员获取所有预约记录
// @Author John 2023-04-28 14:53:50
// @Return reserveVos
// @Return lErr
func (r *ReserveService) GetAllReserveRecords() (reserveVos []vo.ReserveVo, lErr *common.LError) {
	reserveRepository := repository.NewReserveRepository()
	reserveVos, err := reserveRepository.GetAllReserveRecords()
	if err != nil {
		return reserveVos, &common.LError{
			HttpCode: http.StatusOK,
			Msg:      "请求失败",
			Err:      errors.New("获取所有预约记录失败"),
		}
	}
	return reserveVos, nil
}

func NewReserveService() ReserveService {
	return ReserveService{
		DB: common.GetDB(),
	}
}
