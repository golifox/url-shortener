package utils

import (
	"math/rand/v2"
	"net/http"
	"testing"
)

const (
	ip     = "221.248.132.138"
	port   = "8080"
	ipAddr = ip + ":" + port
)

func TestRandomString(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		length := rand.IntN(100)
		result := RandomString(length)
		if len(result) != length {
			t.Error("RandomString() failed")
		}
	})
}

func TestRandomRequestID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		result := RandomRequestID()
		if len(result) != 36 {
			t.Error("RandomRequestID() failed")
		}
	})
}

func TestGetClientIP(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: "",
		}
		result := GetClientIP(r)

		if result != "unknown IP" {
			t.Error("GetClientIP failed")
		}
	})

	t.Run("Success", func(t *testing.T) {
		r := &http.Request{
			RemoteAddr: ipAddr,
		}

		result := GetClientIP(r)

		if result != ip {
			t.Error("GetClientIP failed")
		}
	})
}
