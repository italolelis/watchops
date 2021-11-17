package kinesis

import (
	"errors"

	"github.com/go-redis/redis/v8"
	consumer "github.com/harlow/kinesis-consumer"
	store "github.com/harlow/kinesis-consumer/store/memory"
	"github.com/harlow/kinesis-consumer/store/mysql"
	"github.com/harlow/kinesis-consumer/store/postgres"
	redisStore "github.com/harlow/kinesis-consumer/store/redis"
)

type StoreConfig struct {
	Driver  string `default:"memory"`
	AppName string `split_words:"true" default:"watchops_consumer"`
	Redis   struct {
		Address  string
		Password string
		Username string
		DB       int
	}
	Postgres struct {
		TableName string `split_words:"true"`
		DSN       string
	}
	MySQL struct {
		TableName string `split_words:"true"`
		DSN       string
	}
}

func BuildStore(cfg StoreConfig) (consumer.Store, error) {
	switch cfg.Driver {
	case "redis":
		c := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Address,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		return redisStore.New(cfg.AppName, redisStore.WithClient(c))
	case "postgres":
		return postgres.New(cfg.AppName, cfg.Postgres.TableName, cfg.Postgres.DSN)
	case "mysql":
		return mysql.New(cfg.AppName, cfg.MySQL.TableName, cfg.MySQL.DSN)
	case "memory":
		return store.New(), nil
	default:
		return nil, errors.New("storage driver not supported")
	}
}
