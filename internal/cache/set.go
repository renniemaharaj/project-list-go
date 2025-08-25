package cache

import (
	"encoding/json"
	"time"
)

// setItem stores a value in cache
func setItem[T any](key string, value T, ttl time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	client.Set(ctx, key, data, ttl)
}
