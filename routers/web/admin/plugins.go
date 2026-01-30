// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"fmt"
	"net/http"

	plugin_model "code.gitea.io/gitea/models/plugin"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/context"
	plugin_service "code.gitea.io/gitea/services/plugin"
)

const (
	tplPluginsList   = base.TplName("admin/plugins/list")
	tplPluginsMarket = base.TplName("admin/plugins/market")
)

// PluginsList 插件列表页面
func PluginsList(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.plugins.title")
	ctx.Data["PageIsAdminPlugins"] = true

	manager := plugin_service.GetManager()
	if manager == nil {
		ctx.ServerError("GetManager", fmt.Errorf("plugin manager not initialized"))
		return
	}

	// 获取已安装插件
	installed, err := manager.ListInstalled(ctx)
	if err != nil {
		ctx.ServerError("ListInstalled", err)
		return
	}

	ctx.Data["Plugins"] = installed
	ctx.HTML(http.StatusOK, tplPluginsList)
}

// PluginsMarket 插件市场页面
func PluginsMarket(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.plugins.market")
	ctx.Data["PageIsAdminPlugins"] = true

	manager := plugin_service.GetManager()
	if manager == nil {
		ctx.ServerError("GetManager", fmt.Errorf("plugin manager not initialized"))
		return
	}

	// 获取可用插件
	available, err := manager.ListAvailable()
	if err != nil {
		ctx.Flash.Error(ctx.Tr("admin.plugins.market_error") + ": " + err.Error())
		available = make([]*plugin_model.PluginMetadata, 0)
	}

	// 获取已安装插件列表（用于标记）
	installed, err := manager.ListInstalled(ctx)
	if err != nil {
		ctx.ServerError("ListInstalled", err)
		return
	}

	installedMap := make(map[string]bool)
	for _, p := range installed {
		installedMap[p.PluginID] = true
	}

	ctx.Data["AvailablePlugins"] = available
	ctx.Data["InstalledMap"] = installedMap
	ctx.HTML(http.StatusOK, tplPluginsMarket)
}

// PluginInstall 安装插件
func PluginInstall(ctx *context.Context) {
	pluginID := ctx.FormString("plugin_id")
	version := ctx.FormString("version")

	if pluginID == "" {
		ctx.Flash.Error(ctx.Tr("admin.plugins.plugin_id_required"))
		ctx.Redirect("/admin/plugins/market")
		return
	}

	if version == "" {
		version = "latest"
	}

	manager := plugin_service.GetManager()
	if manager == nil {
		ctx.Flash.Error("Plugin manager not initialized")
		ctx.Redirect("/admin/plugins/market")
		return
	}

	if err := manager.Install(ctx, pluginID, version); err != nil {
		if plugin_model.IsErrPluginAlreadyExist(err) {
			ctx.Flash.Error(ctx.Tr("admin.plugins.already_installed"))
		} else {
			ctx.Flash.Error(ctx.Tr("admin.plugins.install_failed") + ": " + err.Error())
		}
		ctx.Redirect("/admin/plugins/market")
		return
	}

	ctx.Flash.Success(ctx.Tr("admin.plugins.install_success"))
	ctx.Redirect("/admin/plugins")
}

// PluginUninstall 卸载插件
func PluginUninstall(ctx *context.Context) {
	pluginID := ctx.Params("id")

	manager := plugin_service.GetManager()
	if manager == nil {
		ctx.Flash.Error("Plugin manager not initialized")
		ctx.Redirect("/admin/plugins")
		return
	}

	if err := manager.Uninstall(ctx, pluginID); err != nil {
		ctx.Flash.Error(ctx.Tr("admin.plugins.uninstall_failed") + ": " + err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("admin.plugins.uninstall_success"))
	}

	ctx.Redirect("/admin/plugins")
}

// PluginToggle 启用/禁用插件
func PluginToggle(ctx *context.Context) {
	pluginID := ctx.Params("id")
	action := ctx.FormString("action")

	manager := plugin_service.GetManager()
	if manager == nil {
		ctx.Flash.Error("Plugin manager not initialized")
		ctx.Redirect("/admin/plugins")
		return
	}

	var err error
	if action == "enable" {
		err = manager.Enable(ctx, pluginID)
	} else if action == "disable" {
		err = manager.Disable(ctx, pluginID)
	} else {
		ctx.Flash.Error(ctx.Tr("admin.plugins.invalid_action"))
		ctx.Redirect("/admin/plugins")
		return
	}

	if err != nil {
		ctx.Flash.Error(ctx.Tr("admin.plugins.toggle_failed") + ": " + err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("admin.plugins.toggle_success"))
	}

	ctx.Redirect("/admin/plugins")
}
