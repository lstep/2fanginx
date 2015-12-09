package database

import (
	"2fanginx/pwMan"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
)

// UserDatabase is the base JSON where users are stored
type UserDatabase struct {
	Users []pwMan.UserInformation `json:"users"`
}

var (
	userDB *UserDatabase
	dbLock = new(sync.RWMutex)
)

// From excellent article by Karl Seguin (http://openmymind.net/Golang-Hot-Configuration-Reload/)
func init() {
	logrus.Info("Loading database of users")
	LoadDB(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			LoadDB(false)
			logrus.Info("Reloaded database after signal call")
		}
	}()
}

// GetDB returns the database
func GetDB() *UserDatabase {
	dbLock.RLock()
	defer dbLock.RUnlock()
	return userDB
}

// LoadDB is a method to load the database into memory
func LoadDB(fail bool) {
	temp := new(UserDatabase)

	file, e := ioutil.ReadFile("./database.json")
	if e != nil {
		logrus.Errorf("Error loading database: %v\n", e)
		if fail {
			os.Exit(1)
		}
	}

	if err := json.Unmarshal(file, temp); err != nil {
		logrus.Error("Error parsing JSON database: ", err)
		if fail {
			os.Exit(1)
		}
	}

	dbLock.Lock()
	userDB = temp
	dbLock.Unlock()
}

// Save saves the JSON db :-)
func (*UserDatabase) Save() *UserDatabase {
	return nil
}
