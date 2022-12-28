package utils

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var client = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestSetIfLarger_IfNewValueThenSet(t *testing.T) {
	res, err := setIfLarger.Run(context.Background(), client, []string{uuid.NewString()}, 2).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), res)
}

func TestSetIfLarger_IfExistingValueIsSmallerThenSet(t *testing.T) {
	uuid := uuid.NewString()
	_, err := client.Set(context.Background(), uuid, int64(2), time.Minute).Result()
	assert.Nil(t, err)
	
	res, err := setIfLarger.Run(context.Background(), client, []string{uuid}, 5).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(5), res)
}

func TestSetIfLarger_IfExistingValueIsLargerThenDoNotSet(t *testing.T) {
	uuid := uuid.NewString()
	_, err := client.Set(context.Background(), uuid, int64(5), time.Minute).Result()
	assert.Nil(t, err)
	
	res, err := setIfLarger.Run(context.Background(), client, []string{uuid}, 2).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(5), res)
}
