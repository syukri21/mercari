package headermatcher

func MatcherOne(key string) (string, bool) {
	switch key {
	case "Connection", "Content-Length":
		return key, false
	default:
		return key, true
	}
}

func MatcherTwo(key string) (string, bool) {
	switch key {
	case "Connection", "Content-Length", "Authorization":
		return "", false
	default:
		return key, true
	}
}
