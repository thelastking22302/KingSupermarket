package redisDB

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/KingSupermarket/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type singletonRedis struct {
	client *redis.Client
}

var (
	once          sync.Once
	instanceRedis = &singletonRedis{}
)

func GetInstanceRedis() *singletonRedis {
	once.Do(func() {
		instanceRedis = &singletonRedis{}
		instanceRedis.init()
	})
	return instanceRedis
}

func (s *singletonRedis) init() {
	log := logger.GetLogger()
	err := godotenv.Load(".env")
	if err != nil {
		log.Errorf("Error loading .env file: %v", err)
	}
	addr := os.Getenv("REDIS_ADDR")
	pwd := os.Getenv("REDIS_PWD")

	s.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       1,
	})
	ping, err := s.client.Ping(context.Background()).Result()
	if err != nil {
		log.Errorf("Error connecting to redis: %v", err)
	}
	log.Infof("ping: %v\n", ping)
}

func (s *singletonRedis) SaveRefreshToken(refreshToken string) error {
	log := logger.GetLogger()
	err := s.client.Set(context.Background(), "refresh", refreshToken, (7*24*60*60)*time.Second).Err()
	if err != nil {
		log.Errorf("error saving refresh token: %v", err)
		return errors.New("saving token faild")
	}
	log.Infof("token saved successfully")
	return nil
}

func (s *singletonRedis) CheckRefreshToken() (bool, error) {
	log := logger.GetLogger()
	val, err := s.client.Get(context.Background(), "refresh").Result()
	if err == redis.Nil {
		log.Errorf("token invalid")
		return false, nil
	} else if err != nil {
		log.Errorf("error getting token: %v", err)
		return false, err
	}
	// Token tồn tại
	log.Infof("token: %v", val)
	return true, nil
}
