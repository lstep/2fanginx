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

package pluginTOTP

import (
	"2fanginx/database"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gebi/scryptauth"
	"github.com/lstep/gototp"
	"github.com/spf13/viper"
)

func checkPassword(username, password string) bool {
	db := database.GetDB()

	hmacKey := []byte(viper.GetString("hmackey"))
	pwhash, err := scryptauth.New(12, hmacKey)
	if err != nil {
		logrus.Error(err)
		return false
	}

	// Find the user
	for _, item := range db.Users {
		if item.Username == username {
			// found !
			pwCost, hash, salt, err := scryptauth.DecodeBase64(item.ScryptPassword)
			if err != nil {
				logrus.Error(err)
				return false
			}

			ok, err := pwhash.Check(pwCost, hash, []byte(password), salt)
			return ok
		}
	}

	logrus.Infof("Username %s not found in the database", username)
	return false
}

// Authentication is the main method: returncode, username, next_url
func Authentication(w http.ResponseWriter, req *http.Request) (int, string, string) {
	// @TODO: Remove w http.ResponseWriter, req *http.Request to make it independent and ease the tests
	udb := database.GetDB()

	req.ParseForm()

	username, password, code := req.Form.Get("username"),
		req.Form.Get("password"),
		req.Form.Get("code")

	next := req.URL.Query().Get("next")
	if next == "" {
		next = "/"
	}

	logrus.Infof("Trying to authenticate %s", username)
	u := udb.FindUser(username)
	if u == nil {
		logrus.Errorf("Username %s not found in database", username)
		return 1, "", next
	}

	otp, err := gototp.New(u.Init2FA)
	if err != nil {
		logrus.Error(err)
		return 1, "", next
	}

	if checkPassword(username, password) &&
		(code == fmt.Sprintf("%06d", otp.FromNow(-1)) ||
			code == fmt.Sprintf("%06d", otp.Now()) ||
			code == fmt.Sprintf("%06d", otp.FromNow(1))) {

		logrus.Infof("Signing cookie for authentified user %s", username)

		return 0, username, next
	}

	logrus.Error("Failed authentication (pass or OTP) for user ", username)
	return 1, "", next
}
