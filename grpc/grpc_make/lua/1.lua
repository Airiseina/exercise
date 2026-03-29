local key = KEYS[1]

local rate = tonumber(ARGV[1])      -- 每秒生成多少令牌
local capacity = tonumber(ARGV[2])  -- 桶容量
local now = tonumber(ARGV[3])       -- 当前时间戳（毫秒）
local requested = tonumber(ARGV[4]) -- 本次请求需要多少令牌

local data = redis.call("HMGET", key, "tokens", "last_time")
local tokens = tonumber(data[1])
local last_time = tonumber(data[2])

if tokens == nil then
    tokens = capacity
end

if last_time == nil then
    last_time = now
end

local delta = math.max(0, now - last_time)
local filled_tokens = math.min(capacity, tokens + (delta * rate / 1000))
local allowed = filled_tokens >= requested

local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end

redis.call("HMSET", key,
    "tokens", new_tokens,
    "last_time", now
)