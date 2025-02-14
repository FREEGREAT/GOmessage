package utils

import (
	"net/http"

	"github.com/mssola/user_agent"
)

func GetDeviceInfo(userAgent string) map[string]string {
	ua := user_agent.New(userAgent)
	browser, version := ua.Browser()

	deviceType := "desktop"
	if ua.Mobile() {
		deviceType = "mobile"
	}
	return map[string]string{
		"browser":    browser,
		"version":    version,
		"os":         ua.OSInfo().Name,
		"deviceType": deviceType,
	}
}

func GetIP(r *http.Request) string {

	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
