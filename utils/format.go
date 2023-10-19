package utils

import "github.com/google/uuid"

func SafeUUIDToString(id *uuid.UUID) string {
	if id != nil {
		return id.String()
	}
	return ""
}
