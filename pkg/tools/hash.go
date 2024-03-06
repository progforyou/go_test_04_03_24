package tools

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"strings"
)

var passwordSalt = "dating-api-password-ohv6sohK"
var hashSalt = "dating-api-session-ohv6sohK"

func HashPassword(password string) string {
	hash := md5.Sum([]byte(password + ":" + passwordSalt))
	return hex.EncodeToString(hash[:])
}

func NormalizeEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	return email
}

func GenerateToken() (string, error) {
	bytes, err := uuid.New().MarshalBinary()
	if err != nil {
		return "", err
	}
	hash := md5.Sum(append(bytes, []byte(hashSalt)...))
	return hex.EncodeToString(hash[:]), nil
}
