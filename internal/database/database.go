package database

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/dotenv-org/godotenvvault/autoload"
	"github.com/frsfahd/go-proxy/config"
	"github.com/redis/rueidis"
)

type Service interface {
	SetString(string, string) error
	GetString(string) (string, error)
	Health() map[string]string
}

type service struct {
	db rueidis.Client
}

func New(configs config.Config) Service {
	log.Println(configs.Redis.Host, configs.Redis.Port)

	fullAddress := fmt.Sprintf("%s:%s", configs.Redis.Host, configs.Redis.Port)

	rdb, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:      []string{fullAddress},
		Username:         configs.Redis.Username,
		Password:         configs.Redis.Password,
		ConnWriteTimeout: 30 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	s := &service{db: rdb}

	return s
}

// Health returns the health status and statistics of the Redis server.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Default is now 5s
	defer cancel()

	stats := make(map[string]string)

	// Check Redis health and populate the stats map
	stats = s.checkRedisHealth(ctx, stats)

	return stats
}

// checkRedisHealth checks the health of the Redis server and adds the relevant statistics to the stats map.
func (s *service) checkRedisHealth(ctx context.Context, stats map[string]string) map[string]string {
	// Ping the Redis server to check its availability.
	pong := s.db.B().Ping().Message("ping").Build()
	res := s.db.Do(ctx, pong)
	// Note: By extracting and simplifying like this, `log.Fatalf(fmt.Sprintf("db down: %v", err))`
	// can be changed into a standard error instead of a fatal error.
	if res.Error() != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", res.Error().Error()))
	}

	// Redis is up
	stats["redis_status"] = "up"
	stats["redis_message"] = "It's healthy"
	stats["redis_ping_response"] = "pong"

	return stats
}

func (s *service) SetString(key string, data string) error {
	cmd := s.db.B().Set().Key(key).Value(data).Ex(time.Hour).Build()
	res := s.db.Do(context.Background(), cmd)
	if res.Error() != nil {
		return res.Error()
	}

	return nil

}

func (s *service) GetString(key string) (string, error) {
	cmd := s.db.B().Get().Key(key).Cache()
	res := s.db.DoCache(context.Background(), cmd, time.Minute)
	// if res.Error() != nil && !res.IsCacheHit() {
	// 	return "", errors.New("no local cache")
	// }
	if res.Error() != nil && res.IsCacheHit() {
		return "", res.Error()
	}
	data, err := res.ToString()
	if err != nil {
		return "", err
	}
	return data, nil
}
