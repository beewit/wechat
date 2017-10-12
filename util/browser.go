package util

import "strings"

func IsWechatBrowser(userAgent string) bool {
	if strings.Contains(strings.ToLower(userAgent), "micromessenger") {
		return true
	}
	return false
}


