/*
 * @Author: catcatproer 1901749687@qq.com
 * @Date: 2026-08-27 18:55:09
 * @LastEditors: catcatproer 1901749687@qq.com
 * @LastEditTime: 2026-09-01 13:27:19
 * @FilePath: \gin-spa-app-framework-back\internal\models\SystemSettingsModel.go
 * @Description:
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package models

import (
	"gin-spa-app-framework-back/common"

	"gorm.io/gorm"
)

type SystemSettings struct {
	*gorm.Model
	// Id       int    `json:"id" gorm:"primaryKey"`
	SetKey   string `json:"set_key"  gorm:"NOT NULL"`
	SetValue string `json:"set_value"  gorm:"NOT NULL"`
}

//func (SystemSettingsModel) TableName() string {
//	return "system_settings"
//}

// 保存设置
func (sys_settings *SystemSettings) UpdateSysSettings(data []SystemSettings) error {
	var res *gorm.DB
	for i := range data {
		res = common.COM_DB.Model(&sys_settings).Where("set_key = ?", data[i].SetKey).Update("set_value", data[i].SetValue)
		if res.Error != nil {
			break
		}
	}
	return res.Error
}

// 获取系统设置
func (sys_settings *SystemSettings) GetSysSettings() ([]SystemSettings, error) {
	var systemSettings []SystemSettings
	res := common.COM_DB.Model(&sys_settings).Find(&systemSettings)
	return systemSettings, res.Error
}

// 获取OneNet设置
func (sys_settings *SystemSettings) GetSysOneNetSettings() ([]SystemSettings, error) {
	var systemSettings []SystemSettings
	res := common.COM_DB.Model(&sys_settings).Where("set_key LIKE ?", "onenet_%").Find(&systemSettings)
	return systemSettings, res.Error
}
