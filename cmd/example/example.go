package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr: "localhost:49153",
	Password: "", // no password set
	DB: 0,  // use default DB
})

var pool = goredis.NewPool(client)

var rs = redsync.New(pool)

func main() {
	mutex := rs.NewMutex("cool-list")
	mutex.LockContext(ctx)
	client.LPush(ctx, "cool-list", "1")

	fmt.Println(client.LRange(ctx, "cool-list", 0, -1))

	mutex.UnlockContext(ctx)
}
