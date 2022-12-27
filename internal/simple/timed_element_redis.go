package simple

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/pkg/utils"
)

type TimedElementRepoRedis struct {
	Client *redis.Client
	SetKey string
}

func (t TimedElementRepoRedis) Add(element TimedElement) {
	t.Client.ZAdd(context.Background(), t.SetKey, &redis.Z{
		Score:  float64(element.Timestamp),
		Member: element.Data,
	})
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
