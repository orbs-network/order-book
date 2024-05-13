package redisrepo

import (
	"context"
)

func (r *redisRepository) EnumSubKeysOf(ctx context.Context, key string) ([]string, error) {

	var keys []string
	var cursor uint64 = 0
	for {
		// Scan for keys with the given prefix
		result, nextCursor, err := r.client.Scan(ctx, cursor, key+"*", -1).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, result...)

		// If nextCursor is 0, iteration is complete
		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	return keys, nil
}

func (r *redisRepository) ReadStrKey(ctx context.Context, key string) (string, error) {
	return "", nil
}
func (r *redisRepository) WriteStrKey(ctx context.Context, key, val string) error {
	return nil
}
