package level

import "github.com/go-redis/redis"

type LegacyRedis struct {
	client *redis.Client
}

func newLegacyRedis(client *redis.Client) *LegacyRedis {
	return &LegacyRedis{client}
}

// Get - Get value from the hash-map
func (lr *LegacyRedis) Get(ident string, key string) string {
	stringCmd := lr.client.HGet(ident, key)
	return stringCmd.Val()
}

// HKeys - Get a list of the keys in the hash-map
func (lr *LegacyRedis) HKeys(ident string) []string {
	stringSliceCmd := lr.client.HKeys(ident)
	return stringSliceCmd.Val()
}

// Set - Set a value in the hash-map
func (lr *LegacyRedis) Set(ident string, key string, value string) error {
	statusCmd := lr.client.HSet(ident, key, value)
	return statusCmd.Err()
}

// SetM - runs HMSET
func (lr *LegacyRedis) SetM(ident string, set map[string]interface{}) error {
	statusCmd := lr.client.HMSet(ident, set)
	return statusCmd.Err()
}

// Delete - Deletes this key
func (lr *LegacyRedis) Delete(ident string) error {
	statusCmd := lr.client.Del(ident)
	return statusCmd.Err()
}
