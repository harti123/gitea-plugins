// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package plugin

import "fmt"

// ErrPluginNotExist 插件不存在错误
type ErrPluginNotExist struct {
	PluginID string
}

func (err ErrPluginNotExist) Error() string {
	return fmt.Sprintf("plugin does not exist [plugin_id: %s]", err.PluginID)
}

// IsErrPluginNotExist 检查是否为插件不存在错误
func IsErrPluginNotExist(err error) bool {
	_, ok := err.(ErrPluginNotExist)
	return ok
}

// ErrPluginAlreadyExist 插件已存在错误
type ErrPluginAlreadyExist struct {
	PluginID string
}

func (err ErrPluginAlreadyExist) Error() string {
	return fmt.Sprintf("plugin already exists [plugin_id: %s]", err.PluginID)
}

// IsErrPluginAlreadyExist 检查是否为插件已存在错误
func IsErrPluginAlreadyExist(err error) bool {
	_, ok := err.(ErrPluginAlreadyExist)
	return ok
}
