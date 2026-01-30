// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"net/http"

	"code.gitea.io/gitea/models/plugin"
	"github.com/go-chi/chi/v5"
)

// Plugin 插件实例（必须导出）
var Plugin LicenseManagerPlugin

// LicenseManagerPlugin 授权管理插件
type LicenseManagerPlugin struct {
	config map[string]interface{}
}

// Info 获取插件信息
func (p *LicenseManagerPlugin) Info() *plugin.PluginInfo {
	return &plugin.PluginInfo{
		ID:          "license-manager",
		Name:        "授权管理",
		Version:     "1.0.0",
		Description: "设备授权管理插件",
		Author:      "Kiro Team",
		Homepage:    "https://github.com/gitea/license-manager-plugin",
		License:     "MIT",
		Permissions: []string{
			"database.read",
			"database.write",
			"api.create",
			"ui.modify",
		},
		HasRoutes:    true,
		HasAPI:       true,
		HasModels:    true,
		HasTemplates: true,
	}
}

// Init 初始化插件
func (p *LicenseManagerPlugin) Init() error {
	// 初始化配置
	p.config = map[string]interface{}{
		"max_devices_per_user":     10,
		"default_expiry_days":      365,
		"allow_permanent_license":  true,
	}

	// 初始化数据库表
	// TODO: 注册数据库模型

	return nil
}

// RegisterRoutes 注册 Web 路由
func (p *LicenseManagerPlugin) RegisterRoutes(r chi.Router) {
	r.Route("/user/settings/license", func(r chi.Router) {
		r.Get("/", p.licenseList)
		r.Get("/new", p.licenseNew)
		r.Post("/new", p.licenseCreate)
		r.Get("/{id}/edit", p.licenseEdit)
		r.Post("/{id}/edit", p.licenseUpdate)
		r.Post("/{id}/delete", p.licenseDelete)
		r.Post("/{id}/toggle", p.licenseToggle)
	})
}

// RegisterAPIRoutes 注册 API 路由
func (p *LicenseManagerPlugin) RegisterAPIRoutes(r chi.Router) {
	r.Route("/api/v1/license", func(r chi.Router) {
		r.Post("/verify", p.verifyLicense)
		r.Post("/register", p.registerDevice)
	})

	r.Route("/api/v1/user/license", func(r chi.Router) {
		r.Get("/devices", p.listDevices)
		r.Post("/devices", p.createDevice)
		r.Delete("/devices/{id}", p.deleteDevice)
		r.Post("/devices/toggle", p.toggleDevice)
	})
}

// RegisterModels 注册数据库模型
func (p *LicenseManagerPlugin) RegisterModels() []interface{} {
	return []interface{}{
		// &AuthorizedDevice{},
	}
}

// GetTemplatePath 获取模板路径
func (p *LicenseManagerPlugin) GetTemplatePath() string {
	return "./plugins/installed/license-manager/templates"
}

// GetAssetsPath 获取静态资源路径
func (p *LicenseManagerPlugin) GetAssetsPath() string {
	return "./plugins/installed/license-manager/assets"
}

// GetLocalePath 获取国际化文件路径
func (p *LicenseManagerPlugin) GetLocalePath() string {
	return "./plugins/installed/license-manager/locales"
}

// Enable 启用插件
func (p *LicenseManagerPlugin) Enable() error {
	// 启用插件时的逻辑
	return nil
}

// Disable 禁用插件
func (p *LicenseManagerPlugin) Disable() error {
	// 禁用插件时的逻辑
	return nil
}

// Uninstall 卸载插件
func (p *LicenseManagerPlugin) Uninstall() error {
	// 卸载插件时的清理逻辑
	// 例如：删除数据库表、清理配置等
	return nil
}

// GetConfig 获取配置
func (p *LicenseManagerPlugin) GetConfig() map[string]interface{} {
	return p.config
}

// SetConfig 设置配置
func (p *LicenseManagerPlugin) SetConfig(config map[string]interface{}) error {
	p.config = config
	return nil
}

// 以下是路由处理函数的占位实现

func (p *LicenseManagerPlugin) licenseList(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现授权列表
}

func (p *LicenseManagerPlugin) licenseNew(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现新建授权页面
}

func (p *LicenseManagerPlugin) licenseCreate(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现创建授权
}

func (p *LicenseManagerPlugin) licenseEdit(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现编辑授权页面
}

func (p *LicenseManagerPlugin) licenseUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现更新授权
}

func (p *LicenseManagerPlugin) licenseDelete(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现删除授权
}

func (p *LicenseManagerPlugin) licenseToggle(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现切换授权状态
}

func (p *LicenseManagerPlugin) verifyLicense(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现授权验证
}

func (p *LicenseManagerPlugin) registerDevice(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现设备注册
}

func (p *LicenseManagerPlugin) listDevices(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现设备列表
}

func (p *LicenseManagerPlugin) createDevice(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现创建设备
}

func (p *LicenseManagerPlugin) deleteDevice(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现删除设备
}

func (p *LicenseManagerPlugin) toggleDevice(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现切换设备状态
}
