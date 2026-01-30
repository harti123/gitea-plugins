// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package license

import (
	"fmt"
)

// ErrDeviceNotExist 设备不存在错误
type ErrDeviceNotExist struct {
	ID          int64
	MachineCode string
	LicenseKey  string
}

// IsErrDeviceNotExist 检查是否为设备不存在错误
func IsErrDeviceNotExist(err error) bool {
	_, ok := err.(ErrDeviceNotExist)
	return ok
}

func (err ErrDeviceNotExist) Error() string {
	if err.ID > 0 {
		return fmt.Sprintf("device does not exist [id: %d]", err.ID)
	}
	if err.MachineCode != "" {
		return fmt.Sprintf("device does not exist [machine_code: %s]", err.MachineCode)
	}
	if err.LicenseKey != "" {
		return fmt.Sprintf("device does not exist [license_key: %s]", err.LicenseKey)
	}
	return "device does not exist"
}

// ErrDeviceAlreadyExist 设备已存在错误
type ErrDeviceAlreadyExist struct {
	MachineCode string
	LicenseKey  string
}

// IsErrDeviceAlreadyExist 检查是否为设备已存在错误
func IsErrDeviceAlreadyExist(err error) bool {
	_, ok := err.(ErrDeviceAlreadyExist)
	return ok
}

func (err ErrDeviceAlreadyExist) Error() string {
	if err.MachineCode != "" {
		return fmt.Sprintf("device already exists [machine_code: %s]", err.MachineCode)
	}
	if err.LicenseKey != "" {
		return fmt.Sprintf("device already exists [license_key: %s]", err.LicenseKey)
	}
	return "device already exists"
}
