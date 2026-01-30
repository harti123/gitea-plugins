// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package license

import (
	"context"
	"time"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/builder"
)

// AuthorizedDevice 授权设备模型
type AuthorizedDevice struct {
	ID             int64              `xorm:"pk autoincr"`
	UserID         int64              `xorm:"NOT NULL INDEX"` // 所属用户ID
	DeviceID       string             `xorm:"VARCHAR(64) NOT NULL INDEX"`
	MachineCode    string             `xorm:"VARCHAR(64) NOT NULL INDEX"`
	MachineName    string             `xorm:"VARCHAR(200)"`
	LicenseKey     string             `xorm:"VARCHAR(128) NOT NULL INDEX"`
	IsEnabled      bool               `xorm:"NOT NULL DEFAULT true INDEX"`
	ExpiryDate     timeutil.TimeStamp `xorm:"INDEX"`
	CreatedUnix    timeutil.TimeStamp `xorm:"created"`
	UpdatedUnix    timeutil.TimeStamp `xorm:"updated"`
	LastVerifiedAt timeutil.TimeStamp `xorm:"INDEX"`
	Remarks        string             `xorm:"TEXT"`
}

func init() {
	db.RegisterModel(new(AuthorizedDevice))
}

// TableName 表名
func (d *AuthorizedDevice) TableName() string {
	return "authorized_device"
}

// IsValid 检查授权是否有效
func (d *AuthorizedDevice) IsValid() bool {
	if !d.IsEnabled {
		return false
	}
	if !d.ExpiryDate.IsZero() && d.ExpiryDate.AsTime().Before(time.Now()) {
		return false
	}
	return true
}

// UpdateLastVerified 更新最后验证时间
func (d *AuthorizedDevice) UpdateLastVerified(ctx context.Context) error {
	d.LastVerifiedAt = timeutil.TimeStampNow()
	_, err := db.GetEngine(ctx).ID(d.ID).Cols("last_verified_at").Update(d)
	return err
}

// GetDeviceByUserAndMachineCode 根据用户ID和机器码获取设备
func GetDeviceByUserAndMachineCode(ctx context.Context, userID int64, machineCode string) (*AuthorizedDevice, error) {
	device := &AuthorizedDevice{}
	has, err := db.GetEngine(ctx).Where("user_id = ? AND machine_code = ?", userID, machineCode).Get(device)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrDeviceNotExist{MachineCode: machineCode}
	}
	return device, nil
}

// GetDeviceByUserAndLicenseKey 根据用户ID和授权码获取设备
func GetDeviceByUserAndLicenseKey(ctx context.Context, userID int64, licenseKey string) (*AuthorizedDevice, error) {
	device := &AuthorizedDevice{}
	has, err := db.GetEngine(ctx).Where("user_id = ? AND license_key = ?", userID, licenseKey).Get(device)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrDeviceNotExist{LicenseKey: licenseKey}
	}
	return device, nil
}

// GetDeviceByID 根据ID获取设备
func GetDeviceByID(ctx context.Context, id int64) (*AuthorizedDevice, error) {
	device := &AuthorizedDevice{ID: id}
	has, err := db.GetEngine(ctx).Get(device)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrDeviceNotExist{ID: id}
	}
	return device, nil
}

// CreateDevice 创建设备
func CreateDevice(ctx context.Context, device *AuthorizedDevice) error {
	_, err := db.GetEngine(ctx).Insert(device)
	return err
}

// UpdateDevice 更新设备
func UpdateDevice(ctx context.Context, device *AuthorizedDevice) error {
	_, err := db.GetEngine(ctx).ID(device.ID).AllCols().Update(device)
	return err
}

// DeleteDevice 删除设备
func DeleteDevice(ctx context.Context, id int64) error {
	_, err := db.GetEngine(ctx).ID(id).Delete(&AuthorizedDevice{})
	return err
}

// SearchDevicesOptions 搜索设备选项
type SearchDevicesOptions struct {
	db.ListOptions
	UserID    int64  // 用户ID（必需，用于数据隔离）
	Keyword   string
	IsEnabled *bool
}

func (opts *SearchDevicesOptions) toConds() builder.Cond {
	cond := builder.NewCond()
	
	// 必须按用户ID过滤
	cond = cond.And(builder.Eq{"user_id": opts.UserID})
	
	if opts.Keyword != "" {
		cond = cond.And(builder.Or(
			builder.Like{"machine_code", opts.Keyword},
			builder.Like{"machine_name", opts.Keyword},
			builder.Like{"license_key", opts.Keyword},
		))
	}
	if opts.IsEnabled != nil {
		cond = cond.And(builder.Eq{"is_enabled": *opts.IsEnabled})
	}
	return cond
}

// SearchDevices 搜索设备
func SearchDevices(ctx context.Context, opts *SearchDevicesOptions) ([]*AuthorizedDevice, int64, error) {
	sess := db.GetEngine(ctx).Where(opts.toConds())

	if opts.PageSize > 0 {
		sess = db.SetSessionPagination(sess, opts)
	}

	devices := make([]*AuthorizedDevice, 0, opts.PageSize)
	count, err := sess.FindAndCount(&devices)
	return devices, count, err
}

// CountDevices 统计设备数量
func CountDevices(ctx context.Context, userID int64) (int64, error) {
	return db.GetEngine(ctx).Where("user_id = ?", userID).Count(&AuthorizedDevice{})
}

// GetDeviceByUserAndID 根据用户ID和设备ID获取设备（确保数据隔离）
func GetDeviceByUserAndID(ctx context.Context, userID int64, deviceID int64) (*AuthorizedDevice, error) {
	device := &AuthorizedDevice{}
	has, err := db.GetEngine(ctx).Where("user_id = ? AND id = ?", userID, deviceID).Get(device)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrDeviceNotExist{ID: deviceID}
	}
	return device, nil
}
