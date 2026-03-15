package hash

import (
	"hash/fnv"
)

// Consistent returns a deterministic value 0-99 for a given flag key and user ID.
func Consistent(flagKey, userID string) int {
	h := fnv.New32a()
	h.Write([]byte(flagKey + userID))
	return int(h.Sum32() % 100)
}
