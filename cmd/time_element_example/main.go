package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/internal/simple"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/pkg/utils"
)

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:49153",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var repo = &simple.TimedElementRepoRedis{
	Client: client,
	SetKey: "mykey",
}

func main() {
	repo.Add(simple.TimedElement{Timestamp: 1, Data: uuid.NewString()})

	values, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(utils.Map(values, func(item *simple.TimedElement) string { return item.Data }))
}
