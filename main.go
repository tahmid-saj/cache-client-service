package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Redis struct {
	RedisClient redis.Client
}

func NewRedis() (*Redis, error) {
	var client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if client == nil {
		return nil, errors.New("unable to run redis")
	}

	return &Redis{
		RedisClient: *client,
	}, nil
}

func (controller UserController, dataToCache interface{}) GetUserProfile(context *gin.Context) {
	endpoint := context.Request.URL

	cachedKey := endpoint.String()

	response := map[string]interface{}{
		"data": dataToCache,
	}

	dataEncoded, err := json.Marshal(&response)
	if err != nil {
		log.Print(err)
	}

	cacheErr := controller.RedisClient.Set(cachedKey, dataEncoded, 10 * time.Second).Err()
	if cacheErr != nil {
		return cacheErr
	}

	context.JSON(http.StatusOK, response)
}