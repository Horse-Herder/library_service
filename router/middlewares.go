package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"library_server/auth"
	"library_server/auth/jwt"
	"library_server/response"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var claims *jwt.Claims
		var uid string
		var userName string
		var isAdmin bool
		//info := map[string]interface{}{}
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			// 非登录状态
			response.Response(c, http.StatusForbidden, gin.H{
				"status":     http.StatusForbidden,
				"error_code": 0,
				"msg":        "授权失败, token过期，请重新登录",
			})
			c.Abort()
			return
		} else {
			claims, _, err = auth.Jwt.Parse(token)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					response.Response(c, http.StatusForbidden, gin.H{
						"status":     http.StatusForbidden,
						"error_code": 0,
						"msg":        "您已在其他地方登录，如不是本人，请及时修改密码！",
					})
				} else if err == jwt.ErrDeletedToken {
					response.Response(c, http.StatusForbidden, gin.H{
						"status":     http.StatusForbidden,
						"error_code": 0,
						"msg":        "授权已过期，请重新登录",
					})
				} else {
					response.Response(c, http.StatusForbidden, gin.H{
						"status":     http.StatusForbidden,
						"error_code": 0,
						"msg":        "授权失败",
					})

				}
				c.Abort()
				return
			}

			uid = claims.GetUserId()
			userName = claims.GetUserName()
			isAdmin = claims.GetIsAdmin()

			if uid == "" || userName == "" {
				response.Response(c, http.StatusForbidden, gin.H{
					"status":     http.StatusForbidden,
					"error_code": 0,
					"msg":        "授权失败",
				})
			}
			c.Set("uid", uid)
			c.Set("username", userName)
			c.Set("isAdmin", isAdmin)

			c.Set("userid", claims)
		}
	}
}
