/*
 * @Author: catcatproer 1901749687@qq.com
 * @Date: 2026-08-27 18:55:09
 * @LastEditors: catcatproer 1901749687@qq.com
 * @LastEditTime: 2026-09-01 12:55:18
 * @FilePath: \gin-spa-app-framework-back\internal\services\SysService.go
 * @Description:
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package services

import (
	"gin-spa-app-framework-back/internal/models"
	"gin-spa-app-framework-back/internal/models/requests"
	"gin-spa-app-framework-back/internal/models/system"
	"gin-spa-app-framework-back/utils"
)

type SysService struct {
}

func (ss SysService) GenerateCaptcha() (sc system.SysCaptchaModel, err error) {
	var s_captcha system.SysCaptchaModel

	var cap utils.CaptchaInterfaceV2
	cap = new(utils.CaptchaV2)
	id, captch, err := cap.Generate()
	if err != nil {
		return s_captcha, err
	}

	s_captcha.Id = id
	s_captcha.Captcha = captch
	sc = s_captcha
	return

}

func (ss SysService) SystemSettingsUpdate(req *[]requests.SystemSettingsRequest) error {
	var list []models.SystemSettings
	for _, item := range *req {
		// newId, _ := strconv.Atoi(item.Id)
		list = append(list, models.SystemSettings{
			// Id:       newId,
			SetKey:   item.SetKey,
			SetValue: item.SetValue,
		})
	}

	var systemSettings models.SystemSettings
	return systemSettings.UpdateSysSettings(list)
}

func (ss SysService) GetSystemSettings() ([]models.SystemSettings, error) {
	var systemSettings models.SystemSettings
	return systemSettings.GetSysSettings()
}

func (ss SysService) GetSystemOneNetSettings() ([]models.SystemSettings, error) {
	var systemSettings models.SystemSettings
	return systemSettings.GetSysOneNetSettings()
}
