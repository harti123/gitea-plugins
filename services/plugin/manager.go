// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package plugin

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	plugin_model "code.gitea.io/gitea/models/plugin"
	"code.gitea.io/gitea/modules/log"
)

// PluginManager 插件管理器
type PluginManager struct {
	loader         *PluginLoader
	pluginsDir     string
	marketplaceURL string
}

var globalManager *PluginManager

// InitManager 初始化插件管理器
func InitManager(pluginsDir, marketplaceURL string) error {
	globalManager = &PluginManager{
		loader:         GetLoader(),
		pluginsDir:     pluginsDir,
		marketplaceURL: marketplaceURL,
	}

	// 设置插件目录
	globalManager.loader.SetPluginsDir(pluginsDir)

	// 创建必要的目录
	dirs := []string{
		filepath.Join(pluginsDir, "installed"),
		filepath.Join(pluginsDir, "temp"),
		filepath.Join(pluginsDir, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	log.Info("Plugin manager initialized")
	return nil
}

// GetManager 获取全局插件管理器
func GetManager() *PluginManager {
	return globalManager
}

// Install 安装插件
func (m *PluginManager) Install(ctx context.Context, pluginID, version string) error {
	// 检查插件是否已安装
	existing, err := plugin_model.GetPluginByID(ctx, pluginID)
	if err == nil && existing.IsInstalled {
		return plugin_model.ErrPluginAlreadyExist{PluginID: pluginID}
	}

	// 1. 从插件市场下载
	downloadURL := fmt.Sprintf("%s/plugins/%s/%s/download",
		m.marketplaceURL, pluginID, version)

	zipPath := filepath.Join(m.pluginsDir, "temp", pluginID+".zip")
	if err := m.downloadFile(downloadURL, zipPath); err != nil {
		return fmt.Errorf("download plugin: %w", err)
	}
	defer os.Remove(zipPath)

	// 2. 解压到 installed 目录
	installPath := filepath.Join(m.pluginsDir, "installed", pluginID)
	if err := m.unzip(zipPath, installPath); err != nil {
		return fmt.Errorf("unzip plugin: %w", err)
	}

	// 3. 读取插件元数据
	metadata, err := m.loader.readMetadata(installPath)
	if err != nil {
		return fmt.Errorf("read metadata: %w", err)
	}

	// 4. 保存到数据库
	dbPlugin := &plugin_model.Plugin{
		PluginID:    metadata.ID,
		Name:        metadata.Name,
		Version:     metadata.Version,
		Description: metadata.Description,
		Author:      metadata.Author,
		Homepage:    metadata.Homepage,
		License:     metadata.License,
		IsEnabled:   false,
		IsInstalled: true,
		InstallPath: installPath,
	}

	if err := plugin_model.CreatePlugin(ctx, dbPlugin); err != nil {
		return fmt.Errorf("save to database: %w", err)
	}

	// 5. 加载插件
	if err := m.loader.LoadPlugin(ctx, pluginID); err != nil {
		return fmt.Errorf("load plugin: %w", err)
	}

	log.Info("Plugin installed: %s v%s", metadata.Name, metadata.Version)
	return nil
}

// Uninstall 卸载插件
func (m *PluginManager) Uninstall(ctx context.Context, pluginID string) error {
	// 1. 从数据库获取插件信息
	dbPlugin, err := plugin_model.GetPluginByID(ctx, pluginID)
	if err != nil {
		return err
	}

	// 2. 卸载插件
	if err := m.loader.UnloadPlugin(pluginID); err != nil {
		log.Warn("Failed to unload plugin %s: %v", pluginID, err)
	}

	// 3. 删除插件文件
	if err := os.RemoveAll(dbPlugin.InstallPath); err != nil {
		return fmt.Errorf("remove plugin files: %w", err)
	}

	// 4. 从数据库删除
	if err := plugin_model.DeletePlugin(ctx, pluginID); err != nil {
		return fmt.Errorf("delete from database: %w", err)
	}

	log.Info("Plugin uninstalled: %s", pluginID)
	return nil
}

// Enable 启用插件
func (m *PluginManager) Enable(ctx context.Context, pluginID string) error {
	// 1. 获取插件实例
	pluginInstance, ok := m.loader.GetPlugin(pluginID)
	if !ok {
		return fmt.Errorf("plugin not loaded: %s", pluginID)
	}

	// 2. 调用启用方法
	if err := pluginInstance.Enable(); err != nil {
		return fmt.Errorf("enable plugin: %w", err)
	}

	// 3. 更新数据库
	dbPlugin, err := plugin_model.GetPluginByID(ctx, pluginID)
	if err != nil {
		return err
	}

	dbPlugin.IsEnabled = true
	if err := plugin_model.UpdatePlugin(ctx, dbPlugin); err != nil {
		return fmt.Errorf("update database: %w", err)
	}

	log.Info("Plugin enabled: %s", pluginID)
	return nil
}

// Disable 禁用插件
func (m *PluginManager) Disable(ctx context.Context, pluginID string) error {
	// 1. 获取插件实例
	pluginInstance, ok := m.loader.GetPlugin(pluginID)
	if !ok {
		return fmt.Errorf("plugin not loaded: %s", pluginID)
	}

	// 2. 调用禁用方法
	if err := pluginInstance.Disable(); err != nil {
		return fmt.Errorf("disable plugin: %w", err)
	}

	// 3. 更新数据库
	dbPlugin, err := plugin_model.GetPluginByID(ctx, pluginID)
	if err != nil {
		return err
	}

	dbPlugin.IsEnabled = false
	if err := plugin_model.UpdatePlugin(ctx, dbPlugin); err != nil {
		return fmt.Errorf("update database: %w", err)
	}

	log.Info("Plugin disabled: %s", pluginID)
	return nil
}

// ListInstalled 列出已安装的插件
func (m *PluginManager) ListInstalled(ctx context.Context) ([]*plugin_model.Plugin, error) {
	return plugin_model.ListPlugins(ctx)
}

// ListAvailable 列出可用插件（从市场获取）
func (m *PluginManager) ListAvailable() ([]*plugin_model.PluginMetadata, error) {
	resp, err := http.Get(m.marketplaceURL + "/api/plugins")
	if err != nil {
		return nil, fmt.Errorf("fetch from marketplace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("marketplace returned status %d", resp.StatusCode)
	}

	var plugins []*plugin_model.PluginMetadata
	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return plugins, nil
}

// GetPluginInfo 获取插件信息
func (m *PluginManager) GetPluginInfo(pluginID string) (*plugin_model.PluginInfo, error) {
	pluginInstance, ok := m.loader.GetPlugin(pluginID)
	if !ok {
		return nil, fmt.Errorf("plugin not loaded: %s", pluginID)
	}

	return pluginInstance.Info(), nil
}

// downloadFile 下载文件
func (m *PluginManager) downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// unzip 解压文件
func (m *PluginManager) unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// 检查路径安全性
		if !filepath.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
