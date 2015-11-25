package server

// Original from https://gist.github.com/jebjerg/d1c4a23057d5f35a8157 (jebjerg)
// Change CHOOSE-A-SECRET-YOURSELF and eventually the cookie name 'mycookie' currently'

import (
	"2fanginx/pluginTOTP"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
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
	address := viper.GetString("address")
	cookieMaxAge = time.Duration(viper.GetInt("cookiemaxage")) * time.Hour

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

	logrus.Infof("2FA HTTP layer listening on %s", address)
	logrus.Infof("Domain for cookies is %s", viper.GetString("domain"))
	logrus.Infof("Cookie max age is %s hour(s)", viper.GetString("cookiemaxage"))

	if err := http.ListenAndServe(address, httpRateLimiter.RateLimit(mux)); err != nil {
		logrus.Fatal("Unable to create HTTP layer", err)
	}
}
