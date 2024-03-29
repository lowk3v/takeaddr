package utils

import "golang.org/x/exp/slices"

func IsBlacklist(addr string) bool {
	var blacklist []string = []string{
		"0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", // usdc.e
		"0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d", // usdc
		"0x8B39B70E39Aa811b69365398e0aACe9bee238AEb", // usdc
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
