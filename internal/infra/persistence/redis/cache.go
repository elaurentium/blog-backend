package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient encapsula o cliente Redis e fornece métodos para operações de cache.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient cria e retorna uma nova instância do RedisClient.
func NewRedisClient() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // Ex: "redis:6379"
		Password: os.Getenv("REDIS_PASS"), // Sem senha
		DB:       0,                       // Banco de dados padrão
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// Expõe o client interno
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Get recupera um valor do cache pela chave.
func (r *RedisClient) Get(ctx context.Context, key string, result interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), result); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Set armazena um valor no cache com uma chave e um tempo de expiração.
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	if err := r.client.Set(ctx, key, val, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s with expiration %v: %w", key, expiration, err)
	}

	return nil
}

// Delete remove uma chave do cache.
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists verifica se uma chave existe no cache.
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return exists > 0, nil
}

// Close encerra a conexão com o Redis.
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// TTL verifica o tempo de expiração de uma chave.
func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}
	return ttl, nil
}

// FlushAll limpa todos os dados do Redis.
func (r *RedisClient) FlushAll(ctx context.Context) error {
	if err := r.client.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush all keys: %w", err)
	}
	return nil
}

// ScanKeys escaneia as chaves com um padrão específico.
func (r *RedisClient) ScanKeys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	var cursor uint64
	for {
		var err error
		keysBatch, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys with pattern %s: %w", pattern, err)
		}
		keys = append(keys, keysBatch...)
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return keys, nil
}

