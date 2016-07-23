package misc

func InStr(m map[string]string, key string) bool {
	_, ok := m[key]
	return ok
}

func InInt(m map[string]int, key string) bool {
	_, ok := m[key]
	return ok
}
