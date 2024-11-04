package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"net"
	"net/http"
)

// RandomString генерирует случайную строку указанной длины.
func RandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

// RandomRequestID генерирует случайный идентификатор запроса.
func RandomRequestID() string {
	return uuid.New().String()
}

// GetClientIP возвращает IP-адрес клиента.
func GetClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown IP"
	}
	return ip
}
