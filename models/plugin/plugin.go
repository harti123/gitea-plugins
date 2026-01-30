// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package plugin

import (
	"context"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/modules/timeutil"
)

// Plugin 插件数据库模型
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
	Config      string             `xorm:"TEXT"` // JSON 格式的配置
	CreatedUnix timeutil.TimeStamp `xorm:"created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
}

func init() {
	db.RegisterModel(new(Plugin))
}

// TableName 表名
func (p *Plugin) TableName() string {
	return "plugin"
}

// GetPluginByID 根据插件ID获取插件
func GetPluginByID(ctx context.Context, pluginID string) (*Plugin, error) {
	plugin := &Plugin{PluginID: pluginID}
	has, err := db.GetEngine(ctx).Get(plugin)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrPluginNotExist{PluginID: pluginID}
	}
	return plugin, nil
}

// CreatePlugin 创建插件记录
func CreatePlugin(ctx context.Context, plugin *Plugin) error {
	_, err := db.GetEngine(ctx).Insert(plugin)
	return err
}

// UpdatePlugin 更新插件记录
func UpdatePlugin(ctx context.Context, plugin *Plugin) error {
	_, err := db.GetEngine(ctx).ID(plugin.ID).AllCols().Update(plugin)
	return err
}

// DeletePlugin 删除插件记录
func DeletePlugin(ctx context.Context, pluginID string) error {
	_, err := db.GetEngine(ctx).Where("plugin_id = ?", pluginID).Delete(&Plugin{})
	return err
}

// ListPlugins 列出所有插件
func ListPlugins(ctx context.Context) ([]*Plugin, error) {
	plugins := make([]*Plugin, 0)
	err := db.GetEngine(ctx).Find(&plugins)
	return plugins, err
}

// ListEnabledPlugins 列出已启用的插件
func ListEnabledPlugins(ctx context.Context) ([]*Plugin, error) {
	plugins := make([]*Plugin, 0)
	err := db.GetEngine(ctx).Where("is_enabled = ?", true).Find(&plugins)
	return plugins, err
}
