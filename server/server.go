package server

// Original from https://gist.github.com/jebjerg/d1c4a23057d5f35a8157 (jebjerg)
// Change CHOOSE-A-SECRET-YOURSELF and eventually the cookie name 'mycookie' currently'

import (
	"2fanginx/sha512Crypt"
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/craigmj/gototp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

var cookieMaxAge time.Duration

//const cookieMaxAge = 4 * time.Hour

const secureCookie = true

const TOTPSecretPath = "/var/auth/2fa/%v/.google_authenticator"
const shadowFile = "/var/auth/shadow"

func TOTPSecret(user string) (string, error) {
	if len(user) > 0 {
		authFile, err := os.Open(fmt.Sprintf(TOTPSecretPath, user))
		if err != nil {
			return "", err
		}
		defer authFile.Close()
		scanner := bufio.NewScanner(authFile)
		scanner.Scan()
		secret := scanner.Text()
		if len(secret) >= 16 {
			return secret, nil
		}
	}
	return "", fmt.Errorf("bad user '%v'", user)
}

func checkPassword(username, password string) bool {
	shadow, err := os.Open(shadowFile)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	defer shadow.Close()
	scanner := bufio.NewScanner(shadow)
	for scanner.Scan() {
		shadowParts := strings.SplitN(scanner.Text(), ":", 3)
		shadowUser, shadowHash := shadowParts[0], shadowParts[1]
		if shadowUser == username {
			cryptParts := strings.SplitN(shadowHash, "$", 3)
			id := cryptParts[1]
			if id != "6" {
				fmt.Println("WARN! id not 6, refusing")
				return false
			}
			return sha512Crypt.Verify(password, shadowHash)
		}
	}
	return false
}

func freeCookie(w http.ResponseWriter, req *http.Request) {
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

	http.Redirect(w, req, "/", http.StatusFound)
}

func authenticate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username, password, code := req.Form.Get("username"),
		req.Form.Get("password"),
		req.Form.Get("code")

	// TODO: Cleanup the url
	next := req.URL.Query().Get("next")
	if next == "" {
		next = "/"
	}

	secret, err := TOTPSecret(username)
	if err != nil {
		logrus.Error(err)
		http.Redirect(w, req, next, http.StatusTemporaryRedirect)
		return
	}
	otp, err := gototp.New(secret)
	if err != nil {
		logrus.Error(err)
		http.Redirect(w, req, next, http.StatusTemporaryRedirect)
		return
	}

	if checkPassword(username, password) &&
		(code == fmt.Sprintf("%06d", otp.FromNow(-1)) ||
			code == fmt.Sprintf("%06d", otp.Now()) ||
			code == fmt.Sprintf("%06d", otp.FromNow(1))) {
		signResponse(w, username)
		// If has param 'next', go to it, otherwise '/'
		http.Redirect(w, req, next, http.StatusFound)
		return
	}

	http.Redirect(w, req, next, http.StatusTemporaryRedirect)
	return
}

func signResponse(w http.ResponseWriter, username string) {
	expiration := /*username +*/ fmt.Sprintf("%v", int(time.Now().Unix())+3600)
	mac := hmac.New(sha1.New, []byte("CHOOSE-A-SECRET-YOURSELF"))
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
	mux.HandleFunc("/authenticate/free", freeCookie)
	mux.HandleFunc("/authenticate/verify", authenticate)
	mux.Handle("/", http.StripPrefix("/authenticate/", http.FileServer(http.Dir("./static"))))

	logrus.Infof("2FA HTTP layer listening on %s", address)
	logrus.Infof("Domain for cookies is %s", viper.GetString("domain"))
	logrus.Infof("Cookie max age is %s hour(s)", viper.GetString("cookiemaxage"))

	if err := http.ListenAndServe(address, httpRateLimiter.RateLimit(mux)); err != nil {
		logrus.Fatal("Unable to create HTTP layer", err)
	}
}
