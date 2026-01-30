// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package license

import (
	"net/http"
	"time"

	"code.gitea.io/gitea/modules/context"
	license_service "code.gitea.io/gitea/services/license"
)

// VerifyRequest 验证请求
type VerifyRequest struct {
	MachineCode string `json:"machine_code" binding:"required"`
	LicenseKey  string `json:"license_key" binding:"required"`
}

// VerifyResponse 验证响应
type VerifyResponse struct {
	IsAuthorized bool       `json:"is_authorized"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	Message      string     `json:"message"`
	ServerTime   time.Time  `json:"server_time"`
}

// Verify 验证授权
// @Summary 验证设备授权
// @Description 验证设备的机器码和授权码是否有效（需要登录）
// @Tags license
// @Accept json
// @Produce json
// @Param body body VerifyRequest true "验证请求"
// @Success 200 {object} VerifyResponse
// @Router /license/verify [post]
func Verify(ctx *context.APIContext) {
	// 需要用户登录
	if ctx.Doer == nil {
		ctx.Error(http.StatusUnauthorized, "Unauthorized", "需要登录")
		return
	}

	var req VerifyRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Error(http.StatusBadRequest, "bind", err)
		return
	}

	isValid, device, err := license_service.VerifyLicense(ctx, ctx.Doer.ID, req.MachineCode, req.LicenseKey)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "VerifyLicense", err)
		return
	}

	resp := VerifyResponse{
		IsAuthorized: isValid,
		ServerTime:   time.Now(),
	}

	if !isValid {
		if device == nil {
			resp.Message = "设备未授权"
		} else if device.LicenseKey != req.LicenseKey {
			resp.Message = "授权码无效"
		} else if !device.IsEnabled {
			resp.Message = "授权已被禁用"
		} else if !device.ExpiryDate.IsZero() && device.ExpiryDate.AsTime().Before(time.Now()) {
			resp.Message = "授权已过期"
		} else {
			resp.Message = "授权验证失败"
		}
	} else {
		resp.Message = "授权验证成功"
		if !device.ExpiryDate.IsZero() {
			expiryTime := device.ExpiryDate.AsTime()
			resp.ExpiryDate = &expiryTime
		}
	}

	ctx.JSON(http.StatusOK, resp)
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	MachineCode string `json:"machine_code" binding:"required"`
	MachineName string `json:"machine_name"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	MachineCode string `json:"machine_code,omitempty"`
}

// Register 注册设备（仅记录机器码，不生成授权）
// @Summary 注册设备
// @Description 客户端注册设备，记录机器码，等待用户生成授权
// @Tags license
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "注册请求"
// @Success 200 {object} RegisterResponse
// @Router /license/register [post]
func Register(ctx *context.APIContext) {
	// 需要用户登录
	if ctx.Doer == nil {
		ctx.Error(http.StatusUnauthorized, "Unauthorized", "需要登录")
		return
	}

	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Error(http.StatusBadRequest, "bind", err)
		return
	}

	// 检查设备是否已存在
	_, err := license_service.VerifyLicense(ctx, ctx.Doer.ID, req.MachineCode, "")
	if err == nil {
		ctx.JSON(http.StatusOK, RegisterResponse{
			Success: true,
			Message: "设备已注册，请在个人设置中查看授权码",
			MachineCode: req.MachineCode,
		})
		return
	}

	// 返回机器码，用户需要在 Web 界面生成授权
	ctx.JSON(http.StatusOK, RegisterResponse{
		Success: true,
		Message: "请在个人设置的'授权管理'中为此机器码生成授权",
		MachineCode: req.MachineCode,
	})
}
