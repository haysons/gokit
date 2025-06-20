package net

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocalIPs(t *testing.T) {
	ips := GetLocalIPs()

	if ips == nil {
		t.Log("No private ips found or failed to get interfaces (possibly running in a restricted environment)")
	} else {
		for _, ip := range ips {
			t.Log(ip)
			parsedIP := net.ParseIP(ip)
			assert.NotNil(t, parsedIP, "ip should be a valid ip format")
			assert.True(t, parsedIP.IsPrivate(), "ip should be a private address")
		}
	}
}

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()

	if ip == "" {
		t.Log("No private ip found (possibly running in a restricted environment)")
	} else {
		t.Log(ip)
		parsedIP := net.ParseIP(ip)
		assert.NotNil(t, parsedIP, "ip should be a valid ip format")
		assert.True(t, parsedIP.IsPrivate(), "ip should be a private address")
	}
}
