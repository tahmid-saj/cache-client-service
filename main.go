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

type RedisMiddleware struct {
	env         infrastructure.Env
	redisClient infrastructure.Redis
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

func NewRedisMiddleware(
	env infrastructure.Env,
	redisClient infrastructure.Redis,
 ) RedisMiddleware {
	return RedisMiddleware{
	 env:         env,
	 redisClient: redisClient,
	}
 }
 
 // verify Redis Cache
 func (middleware RedisMiddleware) VerifyRedisCache() gin.HandlerFunc {
	return func(context *gin.Context) {
 
	 // get current URL
	 endpoint := context.Request.URL
 
	 // keys are string typed
	 cachedKey := endpoint.String()
 
	 // get Cached keys
	 val, err := middleware.redisClient.RedisClient.Get(cachedKey).Bytes()
 
	 // if error is nil, it means that redis cache couldn't find the key, hence we 
	 // push on to the next middleware to keep the request running
	 if err != nil {
		context.Next()
		return
	 }
 
	 // create an empty interface to unmarshal our cached keys
	 responseBody := map[string]interface{}{}
 
	 // unmarshal cached key
	 json.Unmarshal(val, &responseBody)
 
	 context.JSON(http.StatusOK, responseBody)
	 // abort other chained middlewares since we already get the response here.
	 context.Abort()
	}
 }