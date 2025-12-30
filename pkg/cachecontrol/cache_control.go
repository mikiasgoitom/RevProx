package cachecontrol

import (
	"strconv"
	"strings"
	"time"
)

// type CacheControl struct {
// 	directives map[string]string
// }

func Parse(header string) map[string]string {
	directives := make(map[string]string)
	parts := strings.Split(header, ",")
	for _, part := range parts{
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		key := strings.ToLower(strings.TrimSpace(kv[0]))

		var value string
		if len(kv) == 1 {
			value = ""
		} else {
			value = strings.Trim(strings.TrimSpace(kv[1]), `"`)
		}
		directives[key] = value
		
	}
	return directives
}

func Has(directives map[string]string, directive string) bool {
	_, exists := directives[strings.ToLower(directive)]
	return exists
}

func Get(directives map[string]string, directive string) (string, bool) {
	value, exists := directives[strings.ToLower(directive)]
	return value, exists
}

func GetDuration(directives map[string]string, directive string) (time.Duration, bool) {
	value, exists := directives[strings.ToLower(directive)]
	if !exists {
		return 0, false
	}
	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return time.Duration(seconds) * time.Second, true
}