package utils

import (
	"runtime"
)

func GetOs() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "Mac"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return os
	}
}
