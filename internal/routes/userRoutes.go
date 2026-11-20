/*
 * @Author: catcatproer 1901749687@qq.com
 * @Date: 2026-08-27 18:55:09
 * @LastEditors: catcatproer 1901749687@qq.com
 * @LastEditTime: 2026-09-01 18:38:14
 * @FilePath: \gin-spa-app-framework-back\internal\routes\admin\userRoutes.go
 * @Description: 用户管理路由
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package routes

import (
	controllers "gin-spa-app-framework-back/internal/controllers"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct{}

var UserController = new(controllers.UserController)

func (ur *UserRoutes) InitRouter(PublicRouter *gin.RouterGroup, PrivateRouter *gin.RouterGroup) {
	adminUserGroup := PublicRouter.Group("user")
	{
		adminUserGroup.POST("login", UserController.LoginAction)
		adminUserGroup.POST("test_create", UserController.CreateUserAction)
		adminUserGroup.POST("refresh_token", UserController.RefreshToken)
	}

	adminUserPrivateGroup := PrivateRouter.Group("user")
	{
		adminUserPrivateGroup.POST("create", UserController.CreateUserAction)
		adminUserPrivateGroup.POST("get_user_info", UserController.GetUserInfoAction)
		adminUserPrivateGroup.POST("update", UserController.UpdateUserInfoAction)
	}
}
