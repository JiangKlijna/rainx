package main

import "os"

// File if exist
func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return !stat.IsDir()
	}
	return false
}

// Dir if exist
func DirExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}
	return false
}

const primeRK = 16777619

// Get Ip hash By Address
func IphashByAddress(addr string) uint32 {
	hash := uint32(0)
	n := len(addr) - 1
	for ; n >= 0; n-- {
		if addr[n] == ':' {
			break
		}
	}
	for i := 0; i < n; i++ {
		hash = hash*primeRK + uint32(addr[i])
	}
	return hash
}
