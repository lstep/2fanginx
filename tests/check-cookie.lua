-- set macros: NAME_OF_COOKIE and SIGNING_SECRET_CHOOSE_FOR_YOURSELF
local cookie = "12345454:55fcad87320c34324534285438"
local hmac = ""
local timestamp = ""

local sha1 = require 'sha1'

print("Checking cookie: ", cookie)

function string.tohex(str)
    return (str:gsub('.', function (c)
        return string.format('%02x', string.byte(c))
    end))
end

-- check le cookie
if cookie ~= nil and cookie:find(":") ~= nil then
    -- split cookie into HMAC signature and timestamp.
    local divider = cookie:find(":")
    hmac = cookie:sub(divider+1)
    timestamp = cookie:sub(0, divider-1)
    local secret = "SIGNING_SECRET_CHOOSE_FOR_YOURSELF"

    --local nn = ngx.hmac_sha1(secret, timestamp)
    local nn = sha1.hmac(secret, timestamp)
    if nn ~= hmac then
        print("different hashes")
    end

    print("timestamp",timestamp)
    print("tonumber(timestamp) >= os.time()", tonumber(timestamp) .. ">=" .. os.time())
    print("sub=", tonumber(timestamp) - os.time())
    -- Verify that the signature is valid.
    if tonumber(timestamp) >= os.time() then
       print("OKOK")
       return
    end
end
--
-- redirect no valid cookie found
print("NO valid cookie found")
