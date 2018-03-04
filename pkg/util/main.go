package util

import "log"

// Keys returns the keys from a map
func Keys(structMap map[string]struct{}) []string {
	keys := make([]string, len(structMap))

	i := 0
	for k := range structMap {
		keys[i] = k
		i++
	}
	return keys
}

// Pass returns false if err is not nil and writes err to log
func Pass(err error) bool {
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	return true
}
