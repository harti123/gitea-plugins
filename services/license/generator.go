// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package license

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"code.gitea.io/gitea/modules/util"
)

// GenerateLicenseKey 生成授权码
func GenerateLicenseKey(machineCode string) string {
	data := fmt.Sprintf("%s-%d-%s", machineCode, time.Now().UnixNano(), util.GenerateRandomString(16))
	hash := sha256.Sum256([]byte(data))
	key := hex.EncodeToString(hash[:])[:32]
	return strings.ToUpper(key)
}

// FormatLicenseKey 格式化授权码（XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX）
func FormatLicenseKey(key string) string {
	if len(key) < 32 {
		return key
	}
	parts := make([]string, 0, 8)
	for i := 0; i < 32; i += 4 {
		parts = append(parts, key[i:i+4])
	}
	return strings.Join(parts, "-")
}

// GenerateDeviceID 生成设备ID
func GenerateDeviceID() string {
	return util.GenerateRandomString(16)
}
