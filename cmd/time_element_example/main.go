package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/internal/timedelement"
	"github.com/leonpatmore/go-redis-timestamp-based-generator/pkg/utils"
)

var client = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	Password: "", // no password set
	DB:       0,  // use default DB
})

var repo = timedelement.NewRedisRepo(client, uuid.NewString())

func handleElement(element *timedelement.TimedElement) {
	fmt.Printf("Handling element with ID %s\n", element.Data)
}

func logElements() {
	newValues, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(utils.Map(newValues, func(item *timedelement.TimedElement) string {
		return fmt.Sprintf("%s has timestamp %d", item.Data, item.Timestamp)
	}))
}

func main() {
	fmt.Println("Adding 3 new timed elements...")
	timedelement.AddElementAndHandleIfRequired(repo, &timedelement.TimedElement{Timestamp: 1, Data: uuid.NewString()}, handleElement)
	timedelement.AddElementAndHandleIfRequired(repo, &timedelement.TimedElement{Timestamp: 2, Data: uuid.NewString()}, handleElement)
	timedelement.AddElementAndHandleIfRequired(repo, &timedelement.TimedElement{Timestamp: 3, Data: uuid.NewString()}, handleElement)

	fmt.Println("Before timestamp update:")
	logElements()

	fmt.Println("Pushing timestamp update for timestamp [ 2 ]. This should trigger handling of two elements.")
	timedelement.UpdateTimestampAndHandleElementsBeforeTimestamp(repo, 2, handleElement)

	fmt.Println("After timestamp update:")
	logElements()

	fmt.Println("Pushing element with timestamp [ 2 ]. This element should be handled instantly, and not added to the queue.")
	timedelement.AddElementAndHandleIfRequired(repo, &timedelement.TimedElement{Timestamp: 2, Data: uuid.NewString()}, handleElement)

	fmt.Println("After new element:")
	logElements()

	fmt.Println("Pushing element with timestamp [ 3 ]. This element should not be handled instantly.")
	timedelement.AddElementAndHandleIfRequired(repo, &timedelement.TimedElement{Timestamp: 3, Data: uuid.NewString()}, handleElement)

	fmt.Println("After new element:")
	logElements()

	fmt.Println("Pushing timestamp update for timestamp [ 3 ]. This should trigger handling of two elements.")
	timedelement.UpdateTimestampAndHandleElementsBeforeTimestamp(repo, 3, handleElement)

	fmt.Println("There should now be zero elements left:")
	logElements()
}
