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
	"net/http"
	"time"

	"github.com/lstep/2fanginx/database"
	"github.com/lstep/2fanginx/pluginTOTP"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

var (
	buildst = "none"
	githash = "none"
)

func handleAuthenticate(w http.ResponseWriter, req *http.Request) {

	/* @TODO: Algorithm
	- Get the first in the chain.
	- Execute it, and check its result: 0 -> OK, 1 -> Rejected
	- If there is a next in the chain and the previous result is 0, then continue, otherwise, FINAL Reject.
	- If we are at the last element in the chain, if the result is OK, then FINAL OK, otherwise FINAL Reject.
	*/
	result, username, next := pluginTOTP.Authentication(w, req)
	switch result {
	case 0: // PASSED
		if true { //lastInChain {
			signResponse(w, username)
			http.Redirect(w, req, next, http.StatusFound)
			return
		}

	case 1: // FAILED
		purgeCookie(w)
		http.Redirect(w, req, next, http.StatusTemporaryRedirect)
		return
	}

	// normally should not go here. FAILED
	logrus.Error("Should never arrive here...")
	http.Redirect(w, req, next, http.StatusTemporaryRedirect)
}

// handleFreeCookie sets an invalid/outdated cookie to remove the current one
func handleFreeCookie(w http.ResponseWriter, req *http.Request) {
	// @TODO: get username from the cookie (value returned from purgeCookie?)
	// @TODO: *MUST* store the username in the cookie
	purgeCookie(w)
	logrus.Infof("logged out a user")
	http.Redirect(w, req, "/", http.StatusFound)
}

// Run is the main function
func Run(cmd *cobra.Command, args []string) {
	database.InitDB()
	address := viper.GetString("address")

	// Throttling control
	store, err := memstore.New(65536)
	if err != nil {
		logrus.Fatal(err)
	}

	quota := throttled.RateQuota{throttled.PerMin(20), 5}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		logrus.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/authenticate/free", handleFreeCookie)
	mux.HandleFunc("/authenticate/verify", handleAuthenticate)
	mux.Handle("/", http.StripPrefix("/authenticate/", http.FileServer(http.Dir("./static"))))

	database.InitDB()

	if buildst == "none" {
		buildst = "[This Dev version is not compiled using regular procedure]"
	}

	logrus.Info("Starting App, ", Version)
	logrus.Infof("2FA HTTP layer listening on %s", address)
	logrus.Infof("Domain for cookies is %s", viper.GetString("domain"))
	logrus.Infof("Cookie max age is %d hour(s)", viper.GetInt("cookiemaxage"))
	logrus.Info("Starting instance on ", time.Now())

	if err := http.ListenAndServe(address, httpRateLimiter.RateLimit(mux)); err != nil {
		logrus.Fatal("Unable to create HTTP layer", err)
	}
}
