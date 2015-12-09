package server

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const secureCookie = true

func signResponse(w http.ResponseWriter, username string) {
	// @TODO: Store the username in the cookie (in cleartext) so it can be used afterwards
	cookieMaxAge := time.Duration(viper.GetInt("cookiemaxage")) * time.Hour

	expiration := fmt.Sprintf("%v", int(time.Now().Unix())+int(cookieMaxAge.Seconds()))
	mac := hmac.New(sha1.New, []byte(viper.GetString("cookiesecret")))
	mac.Write([]byte(expiration))
	hash := fmt.Sprintf("%x", mac.Sum(nil))
	value := fmt.Sprintf("%v:%v", expiration, hash)

	cookieContent := fmt.Sprintf("%v=%v", "mycookie", value)
	expire := time.Now().Add(cookieMaxAge)
	cookie := http.Cookie{"mycookie",
		value,
		"/",
		viper.GetString("domain"),
		expire,
		expire.Format(time.UnixDate),
		int(cookieMaxAge.Seconds()),
		secureCookie,
		true,
		cookieContent,
		[]string{cookieContent},
	}
	http.SetCookie(w, &cookie)
}

func purgeCookie(w http.ResponseWriter) {
	cookieContent := fmt.Sprintf("%v=aaaaaa", "mycookie")
	expire := time.Now()
	cookie := http.Cookie{"mycookie",
		"1:1",
		"/",
		viper.GetString("domain"),
		expire,
		expire.Format(time.UnixDate),
		0,
		secureCookie,
		true,
		cookieContent,
		[]string{cookieContent},
	}
	http.SetCookie(w, &cookie)
}
