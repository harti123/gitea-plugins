// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"sync"

	"code.gitea.io/gitea/models/db"
	plugin_model "code.gitea.io/gitea/models/plugin"
	"code.gitea.io/gitea/modules/log"
)

// PluginLoader 插件加载器
type PluginLoader struct {
	pluginsDir string
	plugins    map[string]plugin_model.IPlugin
	mu         sync.RWMutex
}

var (
	globalLoader *PluginLoader
	loaderOnce   sync.Once
)

// GetLoader 获取全局插件加载器
func GetLoader() *PluginLoader {
	loaderOnce.Do(func() {
		globalLoader = &PluginLoader{
			pluginsDir: "./plugins",
			plugins:    make(map[string]plugin.IPlugin),
		}
	})
	return globalLoader
}

// SetPluginsDir 设置插件目录
func (l *PluginLoader) SetPluginsDir(dir string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.pluginsDir = dir
}

// LoadAll 加载所有已安装的插件
func (l *PluginLoader) LoadAll(ctx context.Context) error {
	// 从数据库获取已安装的插件列表
	plugins, err := plugin_model.ListPlugins(ctx)
	if err != nil {
		return fmt.Errorf("list plugins from database: %w", err)
	}

	for _, p := range plugins {
		if !p.IsInstalled {
			continue
		}

		if err := l.LoadPlugin(ctx, p.PluginID); err != nil {
			log.Error("Failed to load plugin %s: %v", p.PluginID, err)
			continue
		}
	}

	log.Info("Loaded %d plugins", len(l.plugins))
	return nil
}

// LoadPlugin 加载单个插件
func (l *PluginLoader) LoadPlugin(ctx context.Context, pluginID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 检查是否已加载
	if _, exists := l.plugins[pluginID]; exists {
		return fmt.Errorf("plugin already loaded: %s", pluginID)
	}

	// 从数据库获取插件信息
	dbPlugin, err := plugin_model.GetPluginByID(ctx, pluginID)
	if err != nil {
		return fmt.Errorf("get plugin from database: %w", err)
	}

	pluginPath := dbPlugin.InstallPath
	if pluginPath == "" {
		pluginPath = filepath.Join(l.pluginsDir, "installed", pluginID)
	}

	// 读取插件元数据
	metadata, err := l.readMetadata(pluginPath)
	if err != nil {
		return fmt.Errorf("read metadata: %w", err)
	}

	// 编译插件（如果需要）
	soPath := filepath.Join(pluginPath, "plugin.so")
	if !fileExists(soPath) {
		log.Info("Compiling plugin: %s", pluginID)
		if err := l.compilePlugin(pluginPath); err != nil {
			return fmt.Errorf("compile plugin: %w", err)
		}
	}

	// 加载插件
	p, err := plugin.Open(soPath)
	if err != nil {
		return fmt.Errorf("open plugin: %w", err)
	}

	// 获取插件实例
	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("lookup Plugin symbol: %w", err)
	}

	pluginInstance, ok := symPlugin.(plugin_model.IPlugin)
	if !ok {
		return fmt.Errorf("invalid plugin type")
	}

	// 初始化插件
	if err := pluginInstance.Init(); err != nil {
		return fmt.Errorf("init plugin: %w", err)
	}

	// 注册数据库模型并创建表
	if metadata.Hooks["models"] {
		models := pluginInstance.RegisterModels()
		if len(models) > 0 {
			if err := l.syncModels(ctx, pluginID, models); err != nil {
				return fmt.Errorf("sync models: %w", err)
			}
			log.Info("Plugin %s synced %d models to database", pluginID, len(models))
		}
	}

	l.plugins[pluginID] = pluginInstance
	log.Info("Plugin loaded: %s v%s", metadata.Name, metadata.Version)

	return nil
}

// syncModels 同步插件的数据库模型
func (l *PluginLoader) syncModels(ctx context.Context, pluginID string, models []interface{}) error {
	if len(models) == 0 {
		return nil
	}

	// 获取数据库引擎
	engine := db.GetEngine(ctx)
	
	// 使用 xorm 的 Sync2 方法创建或更新表结构
	if err := engine.Sync2(models...); err != nil {
		return fmt.Errorf("sync models for plugin %s: %w", pluginID, err)
	}

	log.Info("Successfully synced %d models for plugin %s", len(models), pluginID)
	return nil
}

// UnloadPlugin 卸载插件
func (l *PluginLoader) UnloadPlugin(pluginID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	pluginInstance, exists := l.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not loaded: %s", pluginID)
	}

	// 调用卸载方法
	if err := pluginInstance.Uninstall(); err != nil {
		return fmt.Errorf("uninstall plugin: %w", err)
	}

	delete(l.plugins, pluginID)
	log.Info("Plugin unloaded: %s", pluginID)

	return nil
}

// GetPlugin 获取插件实例
func (l *PluginLoader) GetPlugin(pluginID string) (plugin_model.IPlugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	p, ok := l.plugins[pluginID]
	return p, ok
}

// GetAllPlugins 获取所有已加载的插件
func (l *PluginLoader) GetAllPlugins() map[string]plugin_model.IPlugin {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	result := make(map[string]plugin_model.IPlugin, len(l.plugins))
	for k, v := range l.plugins {
		result[k] = v
	}
	return result
}

// readMetadata 读取插件元数据
func (l *PluginLoader) readMetadata(pluginPath string) (*plugin_model.PluginMetadata, error) {
	metaPath := filepath.Join(pluginPath, "plugin.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read plugin.json: %w", err)
	}

	var metadata plugin_model.PluginMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("parse plugin.json: %w", err)
	}

	return &metadata, nil
}

// compilePlugin 编译插件
func (l *PluginLoader) compilePlugin(pluginPath string) error {
	cmd := exec.Command("go", "build",
		"-buildmode=plugin",
		"-o", "plugin.so",
		".")
	cmd.Dir = pluginPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compile failed: %s\n%s", err, output)
	}

	return nil
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
