/*
 * @Author: catcatproer 1901749687@qq.com
 * @Date: 2026-08-27 18:55:09
 * @LastEditors: catcatproer 1901749687@qq.com
 * @LastEditTime: 2026-08-29 15:13:06
 * @FilePath: \uwb-backed\main.go
 * @Description:
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package main

import (
	"gin-spa-app-framework-back/cmd"
)

func main() {
	//其他服务
	//TODO
	//web服务器
	// res := utils.CheckPasswordHash("123456", "$2a$10$cJGMke/3o8MZ7nDA2LBoDORY2Vk7PfcmSFnKWNw5qWNfKPdyXpjyu")
	// fmt.Println(res)
	cmd.RunServer()
}
