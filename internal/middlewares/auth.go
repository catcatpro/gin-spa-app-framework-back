/*
 * @Author: catcatproer 1901749687@qq.com
 * @Date: 2026-08-27 18:55:09
 * @LastEditors: catcatproer 1901749687@qq.com
 * @LastEditTime: 2026-09-02 10:47:33
 * @FilePath: \gin-spa-app-framework-back\internal\middlewares\auth.go
 * @Description:
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package middlewares

import (
	"fmt"
	"net/http"
	"gin-spa-app-framework-back/utils"

	"github.com/gin-gonic/gin"
)

// 登录鉴权中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-header-token")
		fmt.Println("token", token)
		if token == "" || token == "null" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"msg":    "token not null.",
				"data":   "{}",
			})
			return
		}
		j := utils.NewJWT()
		//验证token是否有效？
		fmt.Println("token: ", token)
		_, err := j.ParseToken(token, false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"msg":    err.Error(),
				"data":   "{}",
			})
			return
		}
		//验证token是否过期？
		_, err = j.VerifyTokenExpiresAt(token, false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"msg":    err.Error(),
				"data":   "{''}",
			})
			return
		}
		c.Next()
	}
}
