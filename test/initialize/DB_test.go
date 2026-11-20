package initialize_test

import (
	"testing"
	"gin-spa-app-framework-back/common"
	"gin-spa-app-framework-back/initialize"
)

func TestConnect(t *testing.T) {
	// configs.InitConfig()
	initialize.InitSystem()
	test_db := common.COM_DB

	if test_db == nil {
		t.Error("db is nil")
	}
}
