if redis.call("EXISTS", KEYS[1]) == 1 then
	return {0, "0"}
end

local amount = tonumber(ARGV[2])
if amount == nil then
	return redis.error_reply("invalid swap amount")
end

local ttl_seconds = tonumber(ARGV[3])
if ttl_seconds == nil or ttl_seconds < 1 then
	return redis.error_reply("invalid swap event ttl")
end

local volume_key_type = redis.call("TYPE", KEYS[2]).ok
if volume_key_type ~= "none" and volume_key_type ~= "hash" then
	return redis.error_reply("invalid swap volume key type")
end

local total_key_type = redis.call("TYPE", KEYS[3]).ok
if total_key_type ~= "none" and total_key_type ~= "string" then
	return redis.error_reply("invalid swap total key type")
end

local current_address_total = redis.call("HGET", KEYS[2], ARGV[1])
if current_address_total ~= false and tonumber(current_address_total) == nil then
	return redis.error_reply("invalid address swap total")
end

local current_total = redis.call("GET", KEYS[3])
if current_total ~= false and tonumber(current_total) == nil then
	return redis.error_reply("invalid swap total")
end

redis.call("SET", KEYS[1], "1", "EX", ttl_seconds)
local address_total = redis.call("HINCRBYFLOAT", KEYS[2], ARGV[1], ARGV[2])
redis.call("INCRBYFLOAT", KEYS[3], ARGV[2])

return {1, address_total}
