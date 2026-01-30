// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

// Plugin settings
var (
	PluginsDir           string
	PluginMarketplaceURL string
	PluginEnabled        bool
	AllowMarketInstall   bool
)

func loadPluginFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("plugin")
	PluginsDir = sec.Key("PLUGINS_DIR").MustString("./plugins")
	PluginMarketplaceURL = sec.Key("MARKETPLACE_URL").MustString("https://plugins.gitea.io")
	PluginEnabled = sec.Key("ENABLED").MustBool(true)
	AllowMarketInstall = sec.Key("ALLOW_MARKETPLACE_INSTALL").MustBool(true)
}
