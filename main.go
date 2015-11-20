// Original from https://gist.github.com/jebjerg/d1c4a23057d5f35a8157 (jebjerg)
// Swap values for  CHANGE FOR YOURSELF, and OBS: it's a novelty authentication, so improvements can and will happen

package main

import (
	"2fanginx/sha512Crypt"
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/craigmj/gototp"
)

// TODO: ?next=..., sha256 cookie

const DOMAIN = ".localhost" // CHANGE FOR YOURSELF
const SECURE_COOKIE = true

const TOTP_SECRET_PATH = "/var/auth/2fa/%v/.google_authenticator"
const SHADOWFILE = "/var/auth/shadow" // CHANGE FOR YOURSELF

func TOTP_Secret(user string) (string, error) {
	if len(user) > 0 {
		auth_file, err := os.Open(fmt.Sprintf(TOTP_SECRET_PATH, user))
		if err != nil {
			return "", err
		}
		defer auth_file.Close()
		scanner := bufio.NewScanner(auth_file)
		scanner.Scan()
		secret := scanner.Text()
		if len(secret) >= 16 {
			return secret, nil
		}
	}
	return "", fmt.Errorf("bad user '%v'", user)
}

func CheckPassword(username, password string) bool {
	shadow, err := os.Open(SHADOWFILE)
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	defer shadow.Close()
	scanner := bufio.NewScanner(shadow)
	for scanner.Scan() {
		shadow_parts := strings.SplitN(scanner.Text(), ":", 3)
		shadow_user, shadow_hash := shadow_parts[0], shadow_parts[1]
		if shadow_user == username {
			crypt_parts := strings.SplitN(shadow_hash, "$", 3)
			id := crypt_parts[1]
			if id != "6" {
				fmt.Println("WARN! id not 6, refusing")
				return false
			}
			return sha512Crypt.Verify(password, shadow_hash)
		}
	}
	return false
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username, password, code := req.Form.Get("username"),
		req.Form.Get("password"),
		req.Form.Get("code")

	secret, err := TOTP_Secret(username)
	if err != nil {
		logrus.Error(err)
		http.Redirect(w, req, "/error", http.StatusTemporaryRedirect)
		return
	}
	otp, err := gototp.New(secret)
	if err != nil {
		logrus.Error(err)
		http.Redirect(w, req, "/error", http.StatusTemporaryRedirect)
		return
	}

	if CheckPassword(username, password) &&
		(code == fmt.Sprintf("%06d", otp.FromNow(-1)) ||
			code == fmt.Sprintf("%06d", otp.Now()) ||
			code == fmt.Sprintf("%06d", otp.FromNow(1))) {
		SignResponse(w, username)
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
	http.Redirect(w, req, "/error", http.StatusTemporaryRedirect)
	return
}

const CookieMaxAge = 4 * time.Hour

func SignResponse(w http.ResponseWriter, username string) {
	expiration := /*username +*/ fmt.Sprintf("%v", int(time.Now().Unix())+3600)
	mac := hmac.New(sha1.New, []byte("cookienamefdjklmgdjfkjfdgkljdf"))
	mac.Write([]byte(expiration))
	hash := fmt.Sprintf("%x", mac.Sum(nil))
	value := fmt.Sprintf("%v:%v", expiration, hash)

	cookieContent := fmt.Sprintf("%v=%v", "cookiename", value)
	expire := time.Now().Add(CookieMaxAge)
	cookie := http.Cookie{"cookiename",
		value,
		"/",
		DOMAIN,
		expire,
		expire.Format(time.UnixDate),
		int(CookieMaxAge.Seconds()),
		SECURE_COOKIE,
		true,
		cookieContent,
		[]string{cookieContent},
	}
	http.SetCookie(w, &cookie)
}

func main() {
	port := ":8080"
	if p := os.Getenv("PORT"); len(p) > 0 {
		port = fmt.Sprintf(":%s", p)
	}

	http.HandleFunc("/authenticate/verify", Authenticate)
	http.Handle("/", http.StripPrefix("/authenticate/", http.FileServer(http.Dir("./static"))))

	fmt.Println("2FA HTTP layer listening")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Unable to create HTTP layer", err)
	}
}
