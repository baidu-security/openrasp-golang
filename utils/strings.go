package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func truncate(s string, n int) string {
	var j int
	for i := range s {
		if j == n {
			return s[:i]
		}
		j++
	}
	return s
}

func TruncateString(s string, n int) string {
	return truncate(s, n)
}

func GetMd5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
