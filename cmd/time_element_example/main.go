package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/internal/simple"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/pkg/utils"
)

var client = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	Password: "", // no password set
	DB:       0,  // use default DB
})

var repo = &simple.TimedElementRepoRedis{
	Client: client,
	SetKey: "mykey",
}

func main() {
	repo.Add(&simple.TimedElement{Timestamp: 1, Data: uuid.NewString()})
	repo.Add(&simple.TimedElement{Timestamp: 2, Data: uuid.NewString()})
	repo.Add(&simple.TimedElement{Timestamp: 3, Data: uuid.NewString()})

	values, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Println("Before:")
	fmt.Println(utils.Map(values, func(item *simple.TimedElement) string { return fmt.Sprintf("%s has score %d", item.Data, item.Timestamp) }))

	err = simple.HandleElementsBeforeTimestamp(repo, 2, func(te *simple.TimedElement) { fmt.Printf("Removing element with ID %s\n", te.Data)})
	if err != nil {
		panic(err)
	}

	fmt.Println("After:")
	newValues, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(utils.Map(newValues, func(item *simple.TimedElement) string { return fmt.Sprintf("%s has score %d", item.Data, item.Timestamp) }))
}
