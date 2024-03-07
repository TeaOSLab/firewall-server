package utils

import (
	"encoding/binary"
	"github.com/cespare/xxhash"
	"math"
	"math/big"
	"net"
	"strings"
)

// IPString2Long 将IP转换为整型
// 注意IPv6没有顺序
func IPString2Long(ip string) uint64 {
	if len(ip) == 0 {
		return 0
	}
	var netIP = net.ParseIP(ip)
	if len(netIP) == 0 {
		return 0
	}
	return NetIP2Long(netIP)
}

// NetIP2Long 将IP对象转换为整型
func NetIP2Long(netIP net.IP) uint64 {
	if len(netIP) == 0 {
		return 0
	}

	var b4 = netIP.To4()
	if b4 != nil {
		return uint64(binary.BigEndian.Uint32(b4.To4()))
	}

	var i = big.NewInt(0)
	i.SetBytes(netIP.To16())
	return i.Uint64()
}

// IP2Long 将IP转换为整型
// 注意IPv6没有顺序
func IP2Long(ip string) uint64 {
	if len(ip) == 0 {
		return 0
	}
	s := net.ParseIP(ip)
	if len(s) == 0 {
		return 0
	}

	if strings.Contains(ip, ":") {
		return math.MaxUint32 + xxhash.Sum64(s)
	}
	return uint64(binary.BigEndian.Uint32(s.To4()))
}

// IsLocalIP 判断是否为本地IP
func IsLocalIP(ipString string) bool {
	var ip = net.ParseIP(ipString)
	if ip == nil {
		return false
	}

	// IPv6
	if strings.Contains(ipString, ":") {
		return ip.String() == "::1"
	}

	// IPv4
	ip = ip.To4()
	if ip == nil {
		return false
	}
	if ip[0] == 127 ||
		ip[0] == 10 ||
		(ip[0] == 172 && ip[1]&0xf0 == 16) ||
		(ip[0] == 192 && ip[1] == 168) {
		return true
	}

	return false
}

// IsIPv4 是否为IPv4
func IsIPv4(ip string) bool {
	var data = net.ParseIP(ip)
	if data == nil {
		return false
	}
	if strings.Contains(ip, ":") {
		return false
	}
	return data.To4() != nil
}

// IsIPv6 是否为IPv6
func IsIPv6(ip string) bool {
	var data = net.ParseIP(ip)
	if data == nil {
		return false
	}
	return !IsIPv4(ip)
}

// IsWildIP 宽泛地判断一个数据是否为IP
func IsWildIP(v string) bool {
	var l = len(v)
	if l == 0 {
		return false
	}

	// for [IPv6]
	if v[0] == '[' && v[l-1] == ']' {
		return IsWildIP(v[1 : l-1])
	}

	return net.ParseIP(v) != nil
}
