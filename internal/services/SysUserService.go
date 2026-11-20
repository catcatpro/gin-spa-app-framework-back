/*
 * @Author: catcatproer 1901749687@qq.com
 * @Date: 2026-08-27 18:55:09
 * @LastEditors: catcatproer 1901749687@qq.com
 * @LastEditTime: 2026-09-01 19:26:27
 * @FilePath: \gin-spa-app-framework-back\internal\services\SysUserService.go
 * @Description: 用户控制器
 *
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved.
 */
package services

import (
	"errors"
	"fmt"
	"gin-spa-app-framework-back/internal/models"
	"gin-spa-app-framework-back/internal/models/requests"
	"gin-spa-app-framework-back/utils"

	"github.com/gin-gonic/gin"
)

type SysUserService struct{}

/**
 * @description:登录
 * @param {*requests.LoginRequest} UserReq
 * @return {*}
 */
func (cs *SysUserService) Login(UserReq *requests.LoginRequest) (string, string, error) {
	var cap utils.CaptchaInterfaceV2
	cap = new(utils.CaptchaV2)
	//验证码验证
	if res, err := cap.Verify(UserReq.CaptchaId, UserReq.Captcha); (err != nil || !res) && gin.Mode() == gin.ReleaseMode {

		return "", "", err
	}

	var userInfo models.User
	userInfo.Username = UserReq.Username
	userInfo.Password = UserReq.Password

	//执行登录验证
	if err := userInfo.UserLogin(); err != nil {
		return "", "", errors.New("Username or Password incorrect")
	}
	// println("debug:", userInfo.ID, userInfo.Username)
	j := utils.NewJWT()

	claims := j.CreateClaims(userInfo.ID, userInfo.Username)
	token, err := j.CreateToken(claims, false)

	if err != nil {
		return "", "", err
	}

	refresh_claims := j.CreateRefreshClaims(userInfo.ID, userInfo.Username)
	refresh_token, err := j.CreateToken(refresh_claims, true)
	if err != nil {
		return "", "", err
	}
	return token, refresh_token, nil
}

/**
 * @description: 刷新token
 * @param {*requests.RefreshTokenRequest} RefreshReq
 * @return {*}
 */
func (cs *SysUserService) RefreshToken(RefreshReq *requests.RefreshTokenRequest) (string, error) {
	j := utils.NewJWT()
	claims, err := j.ParseToken(RefreshReq.RefreshToken, true)
	if err != nil {
		return "", err
	}
	//验证token是否过期？
	_, err = j.VerifyTokenExpiresAt(RefreshReq.RefreshToken, true)
	if err != nil {
		return "", err
	}

	nowClaims := j.CreateClaims(claims.Id, claims.Username)
	now_token, err := j.CreateToken(nowClaims, false)
	if err != nil {
		return "", err
	}
	return now_token, nil
}

/**
 * @description: 新建用户
 * @param {*requests.CreateUserRequest} req
 * @return {*}
 */
func (cs *SysUserService) CreateUser(req *requests.CreateUserRequest) error {

	//确认密码是否正确
	if req.ConfirmPassword != req.Password {
		return errors.New("confirm password not match")
	}

	var userInfo models.User
	userInfo.Username = req.Username
	userInfo.Password = req.Password
	if _, err := userInfo.CreateUser(); err != nil {
		return err
	}

	return nil
}

/**
 * @description: 获取用户信息
 * @param {*requests.GetUserInfoRequest} req
 * @return {*}
 */
func (cs *SysUserService) GetUserInfo(req *requests.GetUserInfoRequest) (models.User, error) {
	var err error = nil
	var userInfo models.User

	switch req.Type {
	case "token":
		err = userInfo.GetUserInfoByToken(req.Data)
		break
	}
	return userInfo, err
}

/**
 * @description: 更新用户信息
 * @param {*requests.UpdateUserInfoRequest} req
 * @return {*}
 */
func (cs *SysUserService) UpdateUserInfo(req *requests.UpdateUserInfoRequest) error {
	var err error = nil
	// basic := &models.Basic{
	// 	ID: req.Id,
	// }
	userInfo := models.User{}
	userInfo.ID = req.Id
	userInfo.Username = req.Username
	userInfo.Password = req.Password
	fmt.Println(userInfo.ID, req.Id)
	err = userInfo.UpdateUserInfo()
	return err
}
