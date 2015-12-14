

[![Total Downloads](https://img.shields.io/github/downloads/2fangnx/latest/total.svg?style=flat-square)](https://github.com/lstep/2fanginx/releases) [![License](http://img.shields.io/badge/license-apache-blue.svg?style=flat-square)](https://raw.githubusercontent.com/lstep/2fanginx/master/LICENSE)


*Documentation is being written right now*

## Purpose

2FANGINX is an auth module for 2FA (2 factors authentication) on NGINX (using "standard" Lua module from NGINX). It allows you to protect using 2FA a whole subdomain, without interfering with other security mesures below the domain hierarchy.

## Features

* Securely hashed (HMAC-SHA1) cookie (distributed only on HTTPS)
* [Throttling connexions](https://github.com/throttled/throttled) to prevent brute force password attempts and DDoS

## References

* Initially based on ([gist](https://gist.github.com/jebjerg/d1c4a23057d5f35a8157) written by [jebjerg](http://github.com/jebjerg))
