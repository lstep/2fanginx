// Original from https://gist.github.com/jebjerg/d1c4a23057d5f35a8157 (jebjerg)
// Change CHOOSE-A-SECRET-YOURSELF and eventually the cookie name 'mycookie' currently'
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

const CookieMaxAge = 4 * time.Hour

const Domain = ".secure.mydomain.eu" // CHANGE FOR YOURSELF
const SecureCookie = true

const TOTPSecretPath = "/var/auth/2fa/%v/.google_authenticator"
const ShadowFile = "/var/auth/shadow"

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

func CheckPassword(username, password string) bool {
	shadow, err := os.Open(ShadowFile)
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
			crypt_parts := strings.SplitN(shadowHash, "$", 3)
			id := crypt_parts[1]
			if id != "6" {
				fmt.Println("WARN! id not 6, refusing")
				return false
			}
			return sha512Crypt.Verify(password, shadowHash)
		}
	}
	return false
}

func FreeCookie(w http.ResponseWriter, req *http.Request) {
	cookieContent := fmt.Sprintf("%v=aaaaaa", "mycookie")
	expire := time.Now()
	cookie := http.Cookie{"mycookie",
		"1:1",
		"/",
		Domain,
		expire,
		expire.Format(time.UnixDate),
		0,
		SecureCookie,
		true,
		cookieContent,
		[]string{cookieContent},
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, req, "/", http.StatusFound)
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
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

	if CheckPassword(username, password) &&
		(code == fmt.Sprintf("%06d", otp.FromNow(-1)) ||
			code == fmt.Sprintf("%06d", otp.Now()) ||
			code == fmt.Sprintf("%06d", otp.FromNow(1))) {
		SignResponse(w, username)
		// If has param 'next', go to it, otherwise '/'
		http.Redirect(w, req, next, http.StatusFound)
		return
	}

	http.Redirect(w, req, next, http.StatusTemporaryRedirect)
	return
}

func SignResponse(w http.ResponseWriter, username string) {
	expiration := /*username +*/ fmt.Sprintf("%v", int(time.Now().Unix())+3600)
	mac := hmac.New(sha1.New, []byte("CHOOSE-A-SECRET-YOURSELF"))
	mac.Write([]byte(expiration))
	hash := fmt.Sprintf("%x", mac.Sum(nil))
	value := fmt.Sprintf("%v:%v", expiration, hash)

	cookieContent := fmt.Sprintf("%v=%v", "mycookie", value)
	expire := time.Now().Add(CookieMaxAge)
	cookie := http.Cookie{"mycookie",
		value,
		"/",
		Domain,
		expire,
		expire.Format(time.UnixDate),
		int(CookieMaxAge.Seconds()),
		SecureCookie,
		true,
		cookieContent,
		[]string{cookieContent},
	}
	http.SetCookie(w, &cookie)
}

func main() {
	port := ":9434"
	if p := os.Getenv("PORT"); len(p) > 0 {
		port = fmt.Sprintf(":%s", p)
	}

	http.HandleFunc("/authenticate/free", FreeCookie)
	http.HandleFunc("/authenticate/verify", Authenticate)
	http.Handle("/", http.StripPrefix("/authenticate/", http.FileServer(http.Dir("./static"))))

	fmt.Println("2FA HTTP layer listening")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Unable to create HTTP layer", err)
	}
}
