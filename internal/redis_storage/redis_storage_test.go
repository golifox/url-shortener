package redis_storage

import (
	"context"
	"testing"
)

const (
	addr        = "localhost:6379"
	invalidAddr = "123.123.1.1:9999"
)

func TestNewRedisStorage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		storage := NewRedisStorage(addr)

		if storage == nil {
			t.Error("Storage is not created")
		}

		err := storage.client.Close()

		if err != nil {
			t.Error("Error closing connection", err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		storage := NewRedisStorage(invalidAddr)

		if storage != nil {
			t.Error("Storage is created")
		}
	})
}

func TestRedisStorage_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		storage := NewRedisStorage(addr)

		ctx := context.Background()

		err := storage.Set(ctx, "key", "value", 0)

		if err != nil {
			t.Error("Error setting value", err)
		}

		_, err = storage.Get(ctx, "key")

		if err != nil {
			t.Error("Error getting value", err)
		}
	})
}

func TestRedisStorage_Set(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		storage := NewRedisStorage(addr)

		ctx := context.Background()

		err := storage.Set(ctx, "key", "value", 0)

		if err != nil {
			t.Error("Error setting value", err)
		}
	})
}

func TestRedisStorage_Exists(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		storage := NewRedisStorage(addr)

		ctx := context.Background()

		err := storage.Set(ctx, "key", "value", 0)

		if err != nil {
			t.Error("Error setting value", err)
		}

		_, err = storage.Exists(ctx, "key")

		if err != nil {
			t.Error("Error checking existence", err)
		}
	})
}

func TestRedisStorage_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		storage := NewRedisStorage(addr)

		ctx := context.Background()

		err := storage.Set(ctx, "key", "value", 0)

		if err != nil {
			t.Error("Error setting value", err)
		}

		err = storage.Delete(ctx, "key")

		if err != nil {
			t.Error("Error deleting value", err)
		}

		_, err = storage.Get(ctx, "key")

		if err == nil {
			t.Error("Value was not deleted")
		}

		if err.Error() != "key not found" {
			t.Error("Unexpected error", err)
		}

	})
}
