/*
In this package we define the API.
*/
package handlers

import (
	"fmt"
	"log"
	"os"
	"redis/data"

	"github.com/go-redis/redis/v8"
)

type API struct {
	l     *log.Logger
	cache *redis.Client
}

func NewApi(l *log.Logger) *API {
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &API{l, rdb}
}

type APIresponse struct {
	Cache bool                     `json:"cache"`
	Data  []data.NominatimResponse `json:"data"`
}
