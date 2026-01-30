// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"
	"strconv"

	"code.gitea.io/gitea/models/license"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/structs"
	license_service "code.gitea.io/gitea/services/license"
)

const (
	tplLicenseList = base.TplName("user/settings/license")
	tplLicenseNew  = base.TplName("user/settings/license_new")
	tplLicenseEdit = base.TplName("user/settings/license_edit")
)

// LicenseList 授权设备列表页面
func LicenseList(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.license.title")
	ctx.Data["PageIsSettingsLicense"] = true

	page := ctx.FormInt("page")
	if page <= 0 {
		page = 1
	}

	opts := &license.SearchDevicesOptions{
		UserID: ctx.Doer.ID, // 只查询当前用户的设备
		ListOptions: structs.ListOptions{
			Page:     page,
			PageSize: setting.UI.User.RepoPagingNum,
		},
	}

	devices, count, err := license.SearchDevices(ctx, opts)
	if err != nil {
		ctx.ServerError("SearchDevices", err)
		return
	}

	ctx.Data["Devices"] = devices
	ctx.Data["Total"] = count

	pager := context.NewPagination(int(count), opts.PageSize, opts.Page, 5)
	ctx.Data["Page"] = pager

	ctx.HTML(http.StatusOK, tplLicenseList)
}

// LicenseNew 新建授权设备页面
func LicenseNew(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.license.new")
	ctx.Data["PageIsSettingsLicense"] = true
	ctx.HTML(http.StatusOK, tplLicenseNew)
}

// LicenseNewPost 处理新建授权设备
func LicenseNewPost(ctx *context.Context) {
	machineCode := ctx.FormString("machine_code")
	machineName := ctx.FormString("machine_name")
	expiryDays := ctx.FormInt("expiry_days")
	remarks := ctx.FormString("remarks")

	if machineCode == "" {
		ctx.Flash.Error(ctx.Tr("settings.license.machine_code_required"))
		ctx.Redirect(setting.AppSubURL + "/user/settings/license/new")
		return
	}

	device, err := license_service.CreateDevice(ctx, &license_service.CreateDeviceOptions{
		UserID:      ctx.Doer.ID, // 当前用户ID
		MachineCode: machineCode,
		MachineName: machineName,
		ExpiryDays:  expiryDays,
		Remarks:     remarks,
	})

	if err != nil {
		if license.IsErrDeviceAlreadyExist(err) {
			ctx.Flash.Error(ctx.Tr("settings.license.device_already_exists"))
		} else {
			ctx.Flash.Error(ctx.Tr("settings.license.create_failed") + ": " + err.Error())
		}
		ctx.Redirect(setting.AppSubURL + "/user/settings/license/new")
		return
	}

	ctx.Flash.Success(ctx.Tr("settings.license.create_success"))
	ctx.Flash.Info("授权码: " + license_service.FormatLicenseKey(device.LicenseKey))
	ctx.Redirect(setting.AppSubURL + "/user/settings/license")
}

// LicenseEdit 编辑授权设备页面
func LicenseEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.license.edit")
	ctx.Data["PageIsSettingsLicense"] = true

	id := ctx.ParamsInt64("id")
	device, err := license.GetDeviceByUserAndID(ctx, ctx.Doer.ID, id)
	if err != nil {
		ctx.NotFound("GetDeviceByUserAndID", err)
		return
	}

	ctx.Data["Device"] = device
	ctx.HTML(http.StatusOK, tplLicenseEdit)
}

// LicenseEditPost 处理编辑授权设备
func LicenseEditPost(ctx *context.Context) {
	id := ctx.ParamsInt64("id")
	machineName := ctx.FormString("machine_name")
	isEnabled := ctx.FormBool("is_enabled")
	expiryDaysStr := ctx.FormString("expiry_days")
	remarks := ctx.FormString("remarks")

	var expiryDays *int
	if expiryDaysStr != "" {
		days, err := strconv.Atoi(expiryDaysStr)
		if err == nil {
			expiryDays = &days
		}
	}

	err := license_service.UpdateDevice(ctx, &license_service.UpdateDeviceOptions{
		UserID:      ctx.Doer.ID, // 当前用户ID
		ID:          id,
		MachineName: machineName,
		IsEnabled:   &isEnabled,
		ExpiryDays:  expiryDays,
		Remarks:     remarks,
	})

	if err != nil {
		ctx.Flash.Error(ctx.Tr("settings.license.update_failed") + ": " + err.Error())
		ctx.Redirect(setting.AppSubURL + "/user/settings/license/" + strconv.FormatInt(id, 10) + "/edit")
		return
	}

	ctx.Flash.Success(ctx.Tr("settings.license.update_success"))
	ctx.Redirect(setting.AppSubURL + "/user/settings/license")
}

// LicenseDelete 删除授权设备
func LicenseDelete(ctx *context.Context) {
	id := ctx.ParamsInt64("id")
	
	// 先检查设备是否属于当前用户
	device, err := license.GetDeviceByUserAndID(ctx, ctx.Doer.ID, id)
	if err != nil {
		ctx.Flash.Error(ctx.Tr("settings.license.device_not_found"))
		ctx.Redirect(setting.AppSubURL + "/user/settings/license")
		return
	}

	if err := license.DeleteDevice(ctx, device.ID); err != nil {
		ctx.Flash.Error(ctx.Tr("settings.license.delete_failed") + ": " + err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("settings.license.delete_success"))
	}
	ctx.Redirect(setting.AppSubURL + "/user/settings/license")
}

// LicenseToggle 切换设备启用状态
func LicenseToggle(ctx *context.Context) {
	id := ctx.ParamsInt64("id")
	if err := license_service.ToggleDevice(ctx, ctx.Doer.ID, id); err != nil {
		ctx.Flash.Error(ctx.Tr("settings.license.toggle_failed") + ": " + err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("settings.license.toggle_success"))
	}
	ctx.Redirect(setting.AppSubURL + "/user/settings/license")
}
