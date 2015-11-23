-- set macros: NAME_OF_COOKIE and SIGNING_SECRET_CHOOSE_FOR_YOURSELF
local cookie = ngx.var.cookie_mycookie
local hmac = ""
local timestamp = ""

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
    local secret = "CHOOSE-A-SECRET-YOURSELF"

    local nn = ngx.hmac_sha1(secret, timestamp)

    -- Verify that the signature is valid.
    if nn:tohex() == hmac and tonumber(timestamp) >= os.time() then
        return
    end
end
--
-- redirect no valid cookie found
ngx.redirect("/authenticate/login.html#next="..ngx.var.uri)
