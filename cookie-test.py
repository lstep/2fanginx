import hmac
from hashlib import sha1

SECRET=b"CHOOSE-COOKIE-SECRET"

h = hmac.new(key=SECRET,msg=b"1448452583", digestmod=sha1)
print(h.hexdigest())

