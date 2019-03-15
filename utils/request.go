package utils

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func GenerateRequestId() string {
	requestId, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return strings.Replace(requestId.String(), "-", "", -1)
}
