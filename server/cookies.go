// Copyright © 2015 Luc Stepniewski <luc@stepniewski.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
