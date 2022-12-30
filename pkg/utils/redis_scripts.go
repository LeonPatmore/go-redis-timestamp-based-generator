package utils

import "github.com/go-redis/redis/v8"

var SetIfLarger = redis.NewScript(`
local key = KEYS[1]
local newValue = tonumber(ARGV[1])
local currentValue = tonumber(redis.call('get', key) or 0)
if (newValue > currentValue) then
    redis.call('set', key, newValue)
    return newValue
end
return currentValue
`)

var AddToSortedSetIfLargerThanNumber = redis.NewScript(`
local setKey = KEYS[1]
local numberKey = KEYS[2]
local targetNumber = tonumber(redis.call('get', numberKey) or 0)
local myScore = tonumber(ARGV[1])
local myData = ARGV[2]
if (myScore > targetNumber) then
    redis.call('zadd', setKey, myScore, myData)
    return myData
end
`)
