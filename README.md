# 2FANGINX

*Documentation is being written right now*

## Purpose

2FANGINX is an auth module for 2FA (2 factors authentication) on NGINX (using "standard" Lua module from NGINX). It allows you to protect using 2FA a whole subdomain, without interfering with other security mesures below the domain hierarchy.

Original ([gist](https://gist.github.com/jebjerg/d1c4a23057d5f35a8157) version was written by [jebjerg](http://github.com/jebjerg)




## Features

* Securely hashed (HMAC-SHA1) cookie (distributed only on HTTPS)
* [Throttling connexions](https://github.com/throttled/throttled) to prevent brute force password attempts and DDoS
