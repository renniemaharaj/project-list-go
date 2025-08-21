package auth

func splitName(name string) (string, string) {
	parts := []rune(name)
	for i, ch := range parts {
		if ch == ' ' {
			return string(parts[:i]), string(parts[i+1:])
		}
	}
	return name, ""
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
