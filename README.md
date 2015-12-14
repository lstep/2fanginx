#2FAnginx

[![Total Downloads](https://img.shields.io/github/downloads/2fangnx/latest/total.svg?style=flat-square)](https://github.com/lstep/2fanginx/releases) [![License](http://img.shields.io/badge/license-apache-blue.svg?style=flat-square)](https://raw.githubusercontent.com/lstep/2fanginx/master/LICENSE) [![Go Report Card](http://goreportcard.com/badge/Masterminds/glide)](http://goreportcard.com/report/lstep/2fanginx)  [![Build Status](https://travis-ci.org/lstep/2fanginx.svg?branch=master)](https://travis-ci.org/lstep/2fanginx)
<span class="badge-patreon"><a href="http://patreon.com/lstep" title="Donate to this project using Patreon"><img src="https://img.shields.io/badge/patreon-donate-yellow.svg" alt="Patreon donate button" /></a></span>
<span class="badge-paypal"><a href="https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=3AZ4NQ7ESWJBC&lc=US&no_note=0&cn=Ajouter%20des%20instructions%20particuli%c3%a8res%20pour%20le%20vendeur%20%3a&no_shipping=2&currency_code=EUR&bn=PP%2dDonationsBF%3abtn_donate_SM%2egif%3aNonHosted" title="Donate to this project using Paypal"><img src="https://img.shields.io/badge/paypal-donate-yellow.svg" alt="PayPal donate button" /></a></span>

*Documentation is being written right now*

## Purpose

2FANGINX is an auth module for 2FA (2 factors authentication) on NGINX (using "standard" Lua module from NGINX). It allows you to protect using 2FA a whole subdomain, without interfering with other security mesures below the domain hierarchy.

## Features

* Securely hashed (HMAC-SHA1) cookie (distributed only on HTTPS)
* [Throttling connexions](https://github.com/throttled/throttled) to prevent brute force password attempts and DDoS

## References

* Initially based on ([gist](https://gist.github.com/jebjerg/d1c4a23057d5f35a8157) written by [jebjerg](http://github.com/jebjerg))
