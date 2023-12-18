package utils

import "golang.org/x/exp/slices"

func IsBlacklist(addr string) bool {
	var blacklist []string = []string{
		"0x0000000000000000000000000000000000000000",
		"0xffffffffffffffffffffffffffffffffffffffff",
		"0xfffffffffffffffffffffffffffffffebaaedce6",
		"0x8000000000000000000000000000000000000000",
		"0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
		"0x7fffffffffffffffffffffffffffffffffffffff",
	}
	if slices.Contains(blacklist, addr) {
		return true
	}
	return false
}
