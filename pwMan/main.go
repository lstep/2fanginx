package pwMan

import (
	"crypto/rand"
	"io"
	"log"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gebi/scryptauth"
	"github.com/spf13/viper"
)

const (
	pwSaltBytes = 32
	pwHashBytes = 64
)

// UserInformation is info about a user
type UserInformation struct {
	Username       string    `json:"username"`
	ScryptPassword string    `json:"scryptpassword"`
	Init2FA        string    `json:"init2fa"`
	LastConnection time.Time `json:"lastconnection"`
	ResetRequired  bool      `json:"resetrequired"`
}

// NewUser creates a user
func NewUser(username string, password string, init2FA string) (*UserInformation, error) {

	hmacKey := []byte(viper.GetString("hmackey"))

	// Create new instace of scryptauth with strength factor 12 and hmac_key
	pwhash, err := scryptauth.New(12, hmacKey)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	hash, salt, err := pwhash.Gen([]byte(password))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	str := scryptauth.EncodeBase64(pwhash.PwCost, hash, salt)
	//fmt.Printf("hash=%x salt=%x\n", hash, salt)
	//fmt.Printf("base64ed: %s\n", str)

	u := &UserInformation{
		Username:       username,
		ScryptPassword: str,
		Init2FA:        init2FA,
		LastConnection: time.Now(),
		ResetRequired:  false,
	}

	return u, nil
}

/*** ***/
func generateSalt() []byte {
	salt := make([]byte, pwSaltBytes)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}
	return salt
}
