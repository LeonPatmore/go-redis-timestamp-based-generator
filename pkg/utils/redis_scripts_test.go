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
	res, err := SetIfLarger.Run(context.Background(), client, []string{uuid.NewString()}, 2).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), res)
}

func TestSetIfLarger_IfExistingValueIsSmallerThenSet(t *testing.T) {
	uuid := uuid.NewString()
	_, err := client.Set(context.Background(), uuid, int64(2), time.Minute).Result()
	assert.Nil(t, err)

	res, err := SetIfLarger.Run(context.Background(), client, []string{uuid}, 5).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(5), res)
}

func TestSetIfLarger_IfExistingValueIsLargerThenDoNotSet(t *testing.T) {
	uuid := uuid.NewString()
	_, err := client.Set(context.Background(), uuid, int64(5), time.Minute).Result()
	assert.Nil(t, err)

	res, err := SetIfLarger.Run(context.Background(), client, []string{uuid}, 2).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(5), res)
}

func TestAddToSortedSetIfLargerThanNumber_NoNumberExists_AddsToSet(t *testing.T) {
	key := uuid.NewString()
	setKey := fmt.Sprintf("set:{%s}", key)
	numberKey := fmt.Sprintf("number:{%s}", key)
	res, err := AddToSortedSetIfLargerThanNumber.Run(context.Background(), client, []string{setKey, numberKey}, 10, "data").Result()

	assert.Nil(t, err)
	assert.Equal(t, "data", res)
}

func TestAddToSortedSetIfLargerThanNumber_NumberExistsAndIsSmaller_AddsToSet(t *testing.T) {
	key := uuid.NewString()
	setKey := fmt.Sprintf("set:{%s}", key)
	numberKey := fmt.Sprintf("number:{%s}", key)

	err := client.Set(context.Background(), numberKey, 6, time.Minute).Err()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestAddToSortedSetIfLargerThanNumber_NumberExistsAndIsSmaller_AddsToSet", func(t *testing.T) {
		res, err := AddToSortedSetIfLargerThanNumber.Run(context.Background(), client, []string{setKey, numberKey}, 10, "data").Result()
		assert.Nil(t, err)
		assert.Equal(t, "data", res)
	})
}

func TestAddToSortedSetIfLargerThanNumber_NumberExistsAndIsLarger_DoNotAddToSet(t *testing.T) {
	key := uuid.NewString()
	setKey := fmt.Sprintf("set:{%s}", key)
	numberKey := fmt.Sprintf("number:{%s}", key)

	err := client.Set(context.Background(), numberKey, 12, time.Minute).Err()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestAddToSortedSetIfLargerThanNumber_NumberExistsAndIsSmaller_AddsToSet", func(t *testing.T) {
		res, err := AddToSortedSetIfLargerThanNumber.Run(context.Background(), client, []string{setKey, numberKey}, 10, "data").Result()
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, res)
	})
}

func TestAddToSortedSetIfLargerThanNumber_NumberExistsAndIsTheSame_DoNotAddToSet(t *testing.T) {
	key := uuid.NewString()
	setKey := fmt.Sprintf("set:{%s}", key)
	numberKey := fmt.Sprintf("number:{%s}", key)

	err := client.Set(context.Background(), numberKey, 10, time.Minute).Err()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestAddToSortedSetIfLargerThanNumber_NumberExistsAndIsSmaller_AddsToSet", func(t *testing.T) {
		res, err := AddToSortedSetIfLargerThanNumber.Run(context.Background(), client, []string{setKey, numberKey}, 10, "data").Result()
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, res)
	})
}
