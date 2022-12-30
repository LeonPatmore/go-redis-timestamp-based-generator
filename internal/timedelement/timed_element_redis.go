package timedelement

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/pkg/utils"
)

type TimedElementRepoRedis struct {
	Client       *redis.Client
	Key          string
	SetKey       string
	TimestampKey string
}

func NewRedisRepo(client *redis.Client, key string) *TimedElementRepoRedis {
	return &TimedElementRepoRedis{
		client,
		key,
		fmt.Sprintf("set:{%s}", key),
		fmt.Sprintf("timestamp:{%s}", key)}
}

func (t TimedElementRepoRedis) AddIfLargerThanTimestamp(element *TimedElement) (bool, error) {
	_, err := utils.AddToSortedSetIfLargerThanNumber.Run(context.Background(), t.Client, []string{t.SetKey, t.TimestampKey}, element.Timestamp, element.Data).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t TimedElementRepoRedis) GetAll() ([]*TimedElement, error) {
	res, err := t.Client.ZRangeWithScores(context.Background(), t.SetKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return utils.Map(res, func(item redis.Z) *TimedElement {
		data := item.Member.(string)
		return &TimedElement{
			Timestamp: int(item.Score),
			Data:      data,
		}
	}), nil
}

func (t TimedElementRepoRedis) Remove(element *TimedElement) error {
	return t.Client.ZRem(context.Background(), t.SetKey, element.Data).Err()
}

func (t TimedElementRepoRedis) GetAllBeforeTimestamp(timestamp int64) ([]*TimedElement, error) {
	res, err := t.Client.ZRangeByScoreWithScores(context.Background(), t.SetKey, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(timestamp, 10),
	}).Result()
	if err != nil {
		return nil, err
	}
	return utils.Map(res, func(item redis.Z) *TimedElement {
		data := item.Member.(string)
		return &TimedElement{
			Timestamp: int(item.Score),
			Data:      data,
		}
	}), nil
}

func (t TimedElementRepoRedis) GetCurrentTimestamp() (*int64, error) {
	res, err := t.Client.Get(context.Background(), t.TimestampKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	resAsInt, err := strconv.ParseInt(res, 10, 64)
	return &resAsInt, err
}

func (t TimedElementRepoRedis) UpdateTimestamp(timestamp int64) (*int64, error) {
	res, err := utils.SetIfLarger.Eval(context.Background(), t.Client, []string{t.TimestampKey}, timestamp).Result()
	resAsInt := res.(int64)
	return &resAsInt, err
}
