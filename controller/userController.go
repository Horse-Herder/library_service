package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	//"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"library_server/auth"
	"library_server/auth/jwt"
	"library_server/common"
	"library_server/model"
	"library_server/response"
	"library_server/service"
	"library_server/utils"
)

type UserController struct {
	DB    *gorm.DB
	Redis common.RedisClient
}

// Register
// @Description 用户注册
// @Author John 2023-04-14 15:22:14
// @Param ctx
func (u *UserController) Register(ctx *gin.Context) {
	//数据接收
	userName := ctx.PostForm("userName")
	email := ctx.PostForm("email")
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")

	reader := model.Reader{
		ReaderName: userName,
		Email:      email,
		Phone:      phone,
		Password:   password,
	}
	userService := service.NewUserService()
	lErr := userService.Register(reader)
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
		"msg":        "注册成功",
	})
}

// Login
// @Description 用户登录
// @Author John 2023-04-14 15:22:25
// @Param ctx
func (u *UserController) Login(ctx *gin.Context) {
	isAdmin := ctx.DefaultPostForm("isAdmin", "false")

	//判断是否管理员
	if isAdmin == "true" {
		u.loginAsAdmin(ctx)
	} else {
		u.loginAsReader(ctx)
	}
}

// loginAsAdmin
// @Description 管理员登陆
// @Author John 2023-04-15 10:57:33
// @Param ctx
func (u *UserController) loginAsAdmin(ctx *gin.Context) {
	// 数据接收
	// 作为管理时，前端接收的数据phone为名称，不需要手机号验证
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")

	admin := model.Admin{
		Phone:    phone,
		Password: password,
	}
	userService := service.NewUserService()
	loginAdmin, lErr := userService.LoginAsAdmin(admin)

	if lErr != nil {
		fmt.Println(lErr.Err)
		response.Response(ctx, lErr.HttpCode, gin.H{
			"status":     lErr.HttpCode,
			"error_code": lErr.ErrorCode,
			"msg":        lErr.Msg,
		})
		return
	}

	ctx.Set("isAdmin", "1")

	claimsElement := &jwt.ClaimsElement{
		UserId:   cast.ToString(loginAdmin.Id),
		UserName: cast.ToString(loginAdmin.Phone),
		IsAdmin:  true,
	}

	token, err := auth.Jwt.Generate(claimsElement, map[string]interface{}{}, 8640)
	if err != nil {
		response.Response(ctx, http.StatusOK, gin.H{
			"status":     http.StatusOK,
			"error_code": 0,
			"msg":        err.Error(),
		})
		return
	}

	cacheKey := fmt.Sprintf("loginAsAdmin_isAdmin%v_token:%s", true, phone)
	u.Redis.Set(ctx, cacheKey, token.AccessToken, 1*time.Hour)

	response.Success(ctx, gin.H{
		"msg":        "管理员登录成功",
		"status":     200,
		"error_code": 1,
		"userName":   phone,
		"isAdmin":    true,
		"token":      token.AccessToken,
	})
}

// loginAsReader
// @Description 读者登录
// @Author John 2023-04-15 10:59:19
// @Param ctx
func (u *UserController) loginAsReader(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")

	// 手机号匹配
	if err := utils.PhoneRegexp(phone); err != nil {
		fmt.Println("请输入正确的手机号")
		response.Response(ctx, http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "请输入正确的手机号",
		})
		return
	}

	loginReader, exist := getReader(u.DB, phone)
	if !exist {
		fmt.Println("账号密码错误或该用户未注册")
		response.Response(ctx, http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "账号密码错误或该用户未注册",
		})
		return
	}
	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(loginReader.Password), []byte(password)); err != nil {
		fmt.Println("账号密码错误或该用户未注册")
		response.Response(ctx, http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "账号密码错误或该用户未注册",
		})
		return
	}

	claimsElement := &jwt.ClaimsElement{
		UserId:   cast.ToString(loginReader.ReaderId),
		UserName: cast.ToString(loginReader.Phone),
		IsAdmin:  false,
	}
	token, err := auth.Jwt.Generate(claimsElement, map[string]interface{}{"isAdmin": false}, 8640)

	if err != nil {
		response.Response(ctx, http.StatusOK, gin.H{
			"status":     http.StatusOK,
			"error_code": 0,
			"msg":        err.Error(),
		})
		return
	}

	cacheKey := fmt.Sprintf("loginAsReader_isAdmin:%v_token:%s", false, loginReader.Phone)
	u.Redis.Set(ctx, cacheKey, token.AccessToken, 1*time.Hour)

	response.Success(ctx, gin.H{
		"msg":         "读者登录成功",
		"status":      200,
		"readerId":    loginReader.ReaderId,
		"readerName":  loginReader.ReaderName,
		"readerPhone": loginReader.Phone,
		"userName":    loginReader.Phone,
		"borrowTimes": loginReader.BorrowTimes,
		"ovdTimes":    loginReader.OvdTimes,
		"email":       loginReader.Email,
		"isAdmin":     false,
		"token":       token.AccessToken,
		"error_code":  1,
	})

}

func (u *UserController) Logout(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	isAdmin := ctx.PostForm("admin")

	cacheKey := fmt.Sprintf("loginAsReader_isAdmin:%v_token:%s", isAdmin, phone)
	u.Redis.Del(ctx, cacheKey)

	response.Success(ctx, gin.H{
		"status":     200,
		"error_code": 1,
		"msg":        "退出登录成功!",
	})
}

func getReader(db *gorm.DB, phone string) (model.Reader, bool) {
	var reader = model.Reader{}
	db.Where("phone = ?", phone).First(&reader)
	return reader, reader.ReaderId != ""
}

// NewUserController
// @Description UserController的构造器
// @Author John 2023-04-16 15:22:31
// @Return UserController
func NewUserController() UserController {
	return UserController{DB: common.GetDB(), Redis: common.GetRedis()}
}
