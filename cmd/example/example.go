package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr: "localhost:49153",
	Password: "", // no password set
	DB: 0,  // use default DB
})

func main() {
	client.LPush(ctx, "cool-list", "1")

	fmt.Println(client.LRange(ctx, "cool-list", 0, -1))
}
