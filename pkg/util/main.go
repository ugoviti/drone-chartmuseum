package util

// Keys returns the keys from a map
func Keys(structMap map[string]bool) []string {
	keys := make([]string, len(structMap))

	i := 0
	for k := range structMap {
		keys[i] = k
		i++
	}
	return keys
}
