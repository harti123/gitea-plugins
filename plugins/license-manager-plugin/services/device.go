// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package license

import (
	"context"
	"time"

	"code.gitea.io/gitea/models/license"
	"code.gitea.io/gitea/modules/timeutil"
)

// CreateDeviceOptions 创建设备选项
type CreateDeviceOptions struct {
	UserID      int64  // 所属用户ID
	MachineCode string
	MachineName string
	ExpiryDays  int // 0 表示永久
	Remarks     string
}

// CreateDevice 创建授权设备
func CreateDevice(ctx context.Context, opts *CreateDeviceOptions) (*license.AuthorizedDevice, error) {
	// 检查该用户的机器码是否已存在
	existing, err := license.GetDeviceByUserAndMachineCode(ctx, opts.UserID, opts.MachineCode)
	if err == nil && existing != nil {
		return nil, license.ErrDeviceAlreadyExist{MachineCode: opts.MachineCode}
	}

	// 生成授权码
	licenseKey := GenerateLicenseKey(opts.MachineCode)

	// 计算到期时间
	var expiryDate timeutil.TimeStamp
	if opts.ExpiryDays > 0 {
		expiryDate = timeutil.TimeStamp(time.Now().AddDate(0, 0, opts.ExpiryDays).Unix())
	}

	device := &license.AuthorizedDevice{
		UserID:      opts.UserID,
		DeviceID:    GenerateDeviceID(),
		MachineCode: opts.MachineCode,
		MachineName: opts.MachineName,
		LicenseKey:  licenseKey,
		IsEnabled:   true,
		ExpiryDate:  expiryDate,
		Remarks:     opts.Remarks,
	}

	if err := license.CreateDevice(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

// VerifyLicense 验证授权
func VerifyLicense(ctx context.Context, userID int64, machineCode, licenseKey string) (bool, *license.AuthorizedDevice, error) {
	device, err := license.GetDeviceByUserAndMachineCode(ctx, userID, machineCode)
	if err != nil {
		if license.IsErrDeviceNotExist(err) {
			return false, nil, nil
		}
		return false, nil, err
	}

	// 验证授权码
	if device.LicenseKey != licenseKey {
		return false, device, nil
	}

	// 检查是否有效
	if !device.IsValid() {
		return false, device, nil
	}

	// 更新最后验证时间
	_ = device.UpdateLastVerified(ctx)

	return true, device, nil
}

// ToggleDevice 切换设备启用状态
func ToggleDevice(ctx context.Context, userID int64, id int64) error {
	device, err := license.GetDeviceByUserAndID(ctx, userID, id)
	if err != nil {
		return err
	}

	device.IsEnabled = !device.IsEnabled
	return license.UpdateDevice(ctx, device)
}

// UpdateDeviceOptions 更新设备选项
type UpdateDeviceOptions struct {
	UserID      int64 // 用户ID（用于数据隔离）
	ID          int64
	MachineName string
	IsEnabled   *bool
	ExpiryDays  *int // nil 表示不修改，0 表示永久，>0 表示天数
	Remarks     string
}

// UpdateDevice 更新设备信息
func UpdateDevice(ctx context.Context, opts *UpdateDeviceOptions) error {
	device, err := license.GetDeviceByUserAndID(ctx, opts.UserID, opts.ID)
	if err != nil {
		return err
	}

	if opts.MachineName != "" {
		device.MachineName = opts.MachineName
	}

	if opts.IsEnabled != nil {
		device.IsEnabled = *opts.IsEnabled
	}

	if opts.ExpiryDays != nil {
		if *opts.ExpiryDays == 0 {
			device.ExpiryDate = timeutil.TimeStamp(0)
		} else {
			device.ExpiryDate = timeutil.TimeStamp(time.Now().AddDate(0, 0, *opts.ExpiryDays).Unix())
		}
	}

	if opts.Remarks != "" {
		device.Remarks = opts.Remarks
	}

	return license.UpdateDevice(ctx, device)
}
