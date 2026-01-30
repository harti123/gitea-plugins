// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_26

import (
	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/xorm"
)

func AddAuthorizedDeviceTable(x *xorm.Engine) error {
	type AuthorizedDevice struct {
		ID          int64              `xorm:"pk autoincr"`
		DeviceID    string             `xorm:"VARCHAR(100) UNIQUE NOT NULL INDEX"`
		UserID      int64              `xorm:"INDEX NOT NULL"`
		MachineCode string             `xorm:"VARCHAR(200) NOT NULL INDEX"`
		MachineName string             `xorm:"VARCHAR(200)"`
		LicenseKey  string             `xorm:"VARCHAR(500) NOT NULL"`
		IsEnabled   bool               `xorm:"NOT NULL DEFAULT true"`
		ExpiryDate  timeutil.TimeStamp `xorm:"INDEX"`
		Remarks     string             `xorm:"TEXT"`
		CreatedUnix timeutil.TimeStamp `xorm:"created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
	}

	return x.Sync2(new(AuthorizedDevice))
}
