-- set macros: NAME_OF_COOKIE and SIGNING_SECRET_CHOOSE_FOR_YOURSELF
local cookie = ngx.var.cookie_NAME_OF_COOKIE
local hmac = ""
local timestamp = ""

-- check le cookie
if cookie ~= nil and cookie:find(":") ~= nil then
    -- split cookie into HMAC signature and timestamp.
    local divider = cookie:find(":")
    hmac = cookie:sub(divider+1)
    timestamp = cookie:sub(0, divider-1)
    local secret = SIGNING_SECRET_CHOOSE_FOR_YOURSELF

    -- Verify that the signature is valid.
    if ndk.set_var.set_encode_hex(ngx.hmac_sha1(secret, timestamp)) == hmac and tonumber(timestamp) >= os.time() then
        return
    end
end
--
-- redirect no valid cookie found
ngx.redirect("/authenticate#next="..ngx.var.uri)