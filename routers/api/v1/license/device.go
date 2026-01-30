// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package license

import (
	"net/http"

	"code.gitea.io/gitea/models/license"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/structs"
	license_service "code.gitea.io/gitea/services/license"
)

// ListDevices 列出当前用户的授权设备
// @Summary 列出授权设备
// @Description 列出当前用户的所有授权设备
// @Tags license
// @Produce json
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {array} license.AuthorizedDevice
// @Router /user/license/devices [get]
func ListDevices(ctx *context.APIContext) {
	// 需要用户登录
	if ctx.Doer == nil {
		ctx.Error(http.StatusUnauthorized, "Unauthorized", "需要登录")
		return
	}

	opts := &license.SearchDevicesOptions{
		UserID: ctx.Doer.ID, // 只查询当前用户的设备
		ListOptions: structs.ListOptions{
			Page:     ctx.FormInt("page"),
			PageSize: ctx.FormInt("limit"),
		},
	}

	if opts.PageSize == 0 {
		opts.PageSize = 20
	}

	devices, count, err := license.SearchDevices(ctx, opts)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "SearchDevices", err)
		return
	}

	ctx.SetTotalCountHeader(count)
	ctx.JSON(http.StatusOK, devices)
}

// CreateDeviceRequest 创建设备请求
type CreateDeviceRequest struct {
	MachineCode string `json:"machine_code" binding:"required"`
	MachineName string `json:"machine_name"`
	ExpiryDays  int    `json:"expiry_days"` // 0 表示永久
	Remarks     string `json:"remarks"`
}

// CreateDeviceResponse 创建设备响应
type CreateDeviceResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	DeviceID   string `json:"device_id,omitempty"`
	LicenseKey string `json:"license_key,omitempty"`
}

// CreateDevice 创建授权设备
// @Summary 创建授权设备
// @Description 为当前用户的指定机器码创建授权并生成授权码
// @Tags license
// @Accept json
// @Produce json
// @Param body body CreateDeviceRequest true "创建请求"
// @Success 200 {object} CreateDeviceResponse
// @Router /user/license/devices [post]
func CreateDevice(ctx *context.APIContext) {
	// 需要用户登录
	if ctx.Doer == nil {
		ctx.Error(http.StatusUnauthorized, "Unauthorized", "需要登录")
		return
	}

	var req CreateDeviceRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Error(http.StatusBadRequest, "bind", err)
		return
	}

	device, err := license_service.CreateDevice(ctx, &license_service.CreateDeviceOptions{
		UserID:      ctx.Doer.ID, // 当前用户ID
		MachineCode: req.MachineCode,
		MachineName: req.MachineName,
		ExpiryDays:  req.ExpiryDays,
		Remarks:     req.Remarks,
	})

	if err != nil {
		if license.IsErrDeviceAlreadyExist(err) {
			ctx.JSON(http.StatusOK, CreateDeviceResponse{
				Success: false,
				Message: "该机器码已存在授权",
			})
			return
		}
		ctx.Error(http.StatusInternalServerError, "CreateDevice", err)
		return
	}

	ctx.JSON(http.StatusOK, CreateDeviceResponse{
		Success:    true,
		Message:    "授权创建成功",
		DeviceID:   device.DeviceID,
		LicenseKey: license_service.FormatLicenseKey(device.LicenseKey),
	})
}

// DeleteDevice 删除授权设备
// @Summary 删除授权设备
// @Description 删除当前用户的指定授权设备
// @Tags license
// @Param id path int true "设备ID"
// @Success 204
// @Router /user/license/devices/{id} [delete]
func DeleteDevice(ctx *context.APIContext) {
	// 需要用户登录
	if ctx.Doer == nil {
		ctx.Error(http.StatusUnauthorized, "Unauthorized", "需要登录")
		return
	}

	id := ctx.ParamsInt64("id")
	
	// 先检查设备是否属于当前用户
	device, err := license.GetDeviceByUserAndID(ctx, ctx.Doer.ID, id)
	if err != nil {
		if license.IsErrDeviceNotExist(err) {
			ctx.Error(http.StatusNotFound, "DeviceNotFound", "设备不存在或无权访问")
			return
		}
		ctx.Error(http.StatusInternalServerError, "GetDevice", err)
		return
	}

	if err := license.DeleteDevice(ctx, device.ID); err != nil {
		ctx.Error(http.StatusInternalServerError, "DeleteDevice", err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// ToggleDeviceRequest 切换设备状态请求
type ToggleDeviceRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// ToggleDevice 切换设备启用状态
// @Summary 切换设备状态
// @Description 启用或禁用当前用户的授权设备
// @Tags license
// @Accept json
// @Produce json
// @Param body body ToggleDeviceRequest true "切换请求"
// @Success 200
// @Router /user/license/devices/toggle [post]
func ToggleDevice(ctx *context.APIContext) {
	// 需要用户登录
	if ctx.Doer == nil {
		ctx.Error(http.StatusUnauthorized, "Unauthorized", "需要登录")
		return
	}

	var req ToggleDeviceRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Error(http.StatusBadRequest, "bind", err)
		return
	}

	if err := license_service.ToggleDevice(ctx, ctx.Doer.ID, req.ID); err != nil {
		if license.IsErrDeviceNotExist(err) {
			ctx.Error(http.StatusNotFound, "DeviceNotFound", "设备不存在或无权访问")
			return
		}
		ctx.Error(http.StatusInternalServerError, "ToggleDevice", err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "设备状态已更新",
	})
}
