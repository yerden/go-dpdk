package lpm

import (
	"encoding/binary"
	"net/netip"
)

func cvtIPv4(addr netip.Addr) (ip uint32) {
	a4 := addr.As4()
	return binary.BigEndian.Uint32(a4[:])
}

func cvtIPv4Net(prefix netip.Prefix) (ip uint32, bits uint8) {
	return cvtIPv4(prefix.Masked().Addr()), uint8(prefix.Bits())
}

func cvtIPv6Net(prefix netip.Prefix) (ip [16]byte, bits uint8) {
	addr := prefix.Masked().Addr()
	if !addr.Is6() || addr.Is4In6() {
		panic("not an IPv6 address")
	}
	return addr.As16(), uint8(prefix.Bits())
}
