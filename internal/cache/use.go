package cache

// Use handles cache-or-fetch logic
func Use[T any](key string, fetch func() (T, error)) (T, error) {
	// 1. Try cache
	if v, found := getItem[T](key); found {
		return v, nil
	}

	// 2. Fetch fresh
	val, err := fetch()
	if err != nil {
		return val, err
	}

	// 3. Store and return
	setItem(key, val, ttl)
	return val, nil
}
