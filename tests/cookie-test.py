import hmac
from hashlib import sha1

SECRET=b"SIGNING_SECRET_CHOOSE_FOR_YOURSELF"
c="123454:343245329954953"
hash=c.split(':')[0]


#h = hmac.new(key=SECRET,msg=b"1448452583", digestmod=sha1)
h = hmac.new(key=SECRET,msg=hash, digestmod=sha1)
print("    Cookie content: %s" % c)
print("Calculated content: %s:%s" % (hash,h.hexdigest()))

