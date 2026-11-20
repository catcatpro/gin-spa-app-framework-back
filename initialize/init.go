package initialize

import (
	"gin-spa-app-framework-back/configs"
)

// 初始化系统
func InitSystem() {
	{
		configs.InitConfig()
		initDB()

	}
}
