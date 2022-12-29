package timedelement

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/pkg/utils"
)

type TimedElementRepoRedis struct {
	Client *redis.Client
	TimestampKey string
	SetKey string
}

func (t TimedElementRepoRedis) Add(element *TimedElement) error {
	return t.Client.ZAdd(context.Background(), t.SetKey, &redis.Z{
		Score:  float64(element.Timestamp),
		Member: element.Data,
	}).Err()
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

func (t TimedElementRepoRedis) GetAllBeforeTimestamp(timestamp int) ([]*TimedElement, error) {
	res, err :=	t.Client.ZRangeByScoreWithScores(context.Background(), t.SetKey, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.Itoa(timestamp),
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

func (t TimedElementRepoRedis) UpdateTimestamp(timestamp int64) error {
	_, err := utils.SetIfLarger.Eval(context.Background(), t.Client, []string{t.TimestampKey}, timestamp).Result()
	return err
}
