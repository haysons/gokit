package net

import (
	"net"
)

// GetLocalIP 获取首个本机私有ip，出现异常或不存在将返回空字符串
func GetLocalIP() string {
	ipAddrs := GetLocalIPs()
	if len(ipAddrs) == 0 {
		return ""
	}
	return ipAddrs[0]
}

// GetLocalIPs 获取本机私有ip列表，出现异常或不存在将返回nil
func GetLocalIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var ipAddrs []string
	for _, iface := range ifaces {
		// 忽略未启用网卡或回环网卡
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// 忽略回环地址、链路本地地址、公网地址
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || !ip.IsPrivate() {
				continue
			}
			ipAddrs = append(ipAddrs, ip.String())
		}
	}
	return ipAddrs
}
