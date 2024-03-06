package tools

import "strings"

func IsErrUniqueConstraint(err error) bool {
	if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") { // SQLITE
		return true
	}
	if strings.HasPrefix(err.Error(), "Error 1062:") { // MySQL
		return true
	}
	return false
}
