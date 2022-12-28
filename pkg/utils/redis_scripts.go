package utils

import "github.com/go-redis/redis/v8"

var setIfLarger = redis.NewScript(`
local key = KEYS[1]
local newValue = tonumber(ARGV[1])
local currentValue = tonumber(redis.call('get', key) or 0)
if (newValue > currentValue) then
    redis.call('set', key, newValue)
    return newValue
end
return currentValue
`)
