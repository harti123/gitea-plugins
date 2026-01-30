// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_26

import (
	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/xorm"
)

func AddPluginTable(x *xorm.Engine) error {
	type Plugin struct {
		ID          int64              `xorm:"pk autoincr"`
		PluginID    string             `xorm:"VARCHAR(100) UNIQUE NOT NULL INDEX"`
		Name        string             `xorm:"VARCHAR(200) NOT NULL"`
		Version     string             `xorm:"VARCHAR(50) NOT NULL"`
		Description string             `xorm:"TEXT"`
		Author      string             `xorm:"VARCHAR(200)"`
		Homepage    string             `xorm:"VARCHAR(500)"`
		License     string             `xorm:"VARCHAR(50)"`
		IsEnabled   bool               `xorm:"NOT NULL DEFAULT false"`
		IsInstalled bool               `xorm:"NOT NULL DEFAULT false"`
		InstallPath string             `xorm:"VARCHAR(500)"`
		Config      string             `xorm:"TEXT"`
		CreatedUnix timeutil.TimeStamp `xorm:"created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
	}

	return x.Sync2(new(Plugin))
}
