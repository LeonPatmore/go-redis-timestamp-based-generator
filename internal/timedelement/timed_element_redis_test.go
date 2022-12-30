package timedelement

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var client = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	Password: "", // no password set
	DB:       0,  // use default DB
})

var repo TimedElementRepo

func setupRandomRepo() {
	repo = *NewRedisRepo(client, uuid.NewString())
}

func runWithRandomRepo(t *testing.T, name string, f func(*testing.T)) bool {
	setupRandomRepo()
	return t.Run(name, f)
}

func TestWithRandomRepo(t *testing.T) {

	runWithRandomRepo(t, "TestUpdateTimestamp_whenNoCurrentTimestamp", func(t *testing.T) {
		res, err := repo.UpdateTimestamp(5)
		assert.Nil(t, err)
		assert.Equal(t, int64(5), *res)
	})

	runWithRandomRepo(t, "TestUpdateTimestamp_whenCurrentTimestampIsLower", func(t *testing.T) {
		_, err := repo.UpdateTimestamp(5)
		assert.Nil(t, err)

		t.Run("TestUpdateTimestamp_whenCurrentTimestampIsLower", func(t *testing.T) {
			res, err := repo.UpdateTimestamp(10)
			assert.Nil(t, err)
			assert.Equal(t, int64(10), *res)
		})
	})

	runWithRandomRepo(t, "TestUpdateTimestamp_whenCurrentTimestampIsHigher", func(t *testing.T) {
		_, err := repo.UpdateTimestamp(10)
		assert.Nil(t, err)

		t.Run("TestUpdateTimestamp_whenCurrentTimestampIsLower", func(t *testing.T) {
			res, err := repo.UpdateTimestamp(6)
			assert.Nil(t, err)
			assert.Equal(t, int64(10), *res)
		})
	})

	runWithRandomRepo(t, "TestGetCurrentTimestamp_WhenNoCurrentTimestamp", func(t *testing.T) {
		res, err := repo.GetCurrentTimestamp()
		assert.Nil(t, err)
		assert.Equal(t, (*int64)(nil), res)
	})

	runWithRandomRepo(t, "TestGetCurrentTimestamp_WhenCurrentTimestampExists", func(t *testing.T) {
		_, err := repo.UpdateTimestamp(5)
		assert.Nil(t, err)
		t.Run("TestGetCurrentTimestamp_WhenCurrentTimestamp", func(t *testing.T) {
			res, err := repo.GetCurrentTimestamp()

			assert.Nil(t, err)
			assert.Equal(t, int64(5), *res)
		})
	})

}
