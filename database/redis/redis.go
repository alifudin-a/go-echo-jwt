package redis

import (
	"os"

	"github.com/go-redis/redis/v7"
	_ "github.com/joho/godotenv/autoload"
)

// RdbClient global variable
var RdbClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_DSN"),
	Password: "",
	DB:       1,
})
