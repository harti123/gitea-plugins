// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package plugin

import (
	"github.com/go-chi/chi/v5"
)

// IPlugin 插件接口
type IPlugin interface {
	// Info 获取插件信息
	Info() *PluginInfo

	// Init 初始化插件
	Init() error

	// RegisterRoutes 注册 Web 路由
	RegisterRoutes(r chi.Router)

	// RegisterAPIRoutes 注册 API 路由
	RegisterAPIRoutes(r chi.Router)

	// RegisterModels 注册数据库模型
	RegisterModels() []interface{}

	// GetTemplatePath 获取模板路径
	GetTemplatePath() string

	// GetAssetsPath 获取静态资源路径
	GetAssetsPath() string

	// GetLocalePath 获取国际化文件路径
	GetLocalePath() string

	// Enable 启用插件
	Enable() error

	// Disable 禁用插件
	Disable() error

	// Uninstall 卸载插件
	Uninstall() error

	// GetConfig 获取配置
	GetConfig() map[string]interface{}

	// SetConfig 设置配置
	SetConfig(config map[string]interface{}) error
}

// PluginInfo 插件信息
type PluginInfo struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Description    string            `json:"description"`
	Author         string            `json:"author"`
	Homepage       string            `json:"homepage"`
	License        string            `json:"license"`
	GiteaVersion   string            `json:"gitea_version"`
	Dependencies   []string          `json:"dependencies"`
	Permissions    []string          `json:"permissions"`
	ConfigSchema   map[string]interface{} `json:"config_schema"`
	HasRoutes      bool              `json:"has_routes"`
	HasAPI         bool              `json:"has_api"`
	HasModels      bool              `json:"has_models"`
	HasTemplates   bool              `json:"has_templates"`
}

// PluginMetadata 插件元数据（从 plugin.json 读取）
type PluginMetadata struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author"`
	Homepage     string                 `json:"homepage"`
	License      string                 `json:"license"`
	GiteaVersion string                 `json:"gitea_version"`
	Dependencies []string               `json:"dependencies"`
	EntryPoint   string                 `json:"entry_point"`
	Hooks        map[string]bool        `json:"hooks"`
	Permissions  []string               `json:"permissions"`
	ConfigSchema map[string]interface{} `json:"config_schema"`
}
