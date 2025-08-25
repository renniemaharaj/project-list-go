package cache

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// getItem tries to fetch a cached value
func getItem[T any](key string) (T, bool) {
	var v T
	raw, err := client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		var zero T
		return zero, false
	}
	if jsonErr := json.Unmarshal([]byte(raw), &v); jsonErr != nil {
		var zero T
		return zero, false
	}
	return v, true
}
