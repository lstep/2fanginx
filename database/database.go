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
	"github.com/spf13/viper"
)

// UserDatabase is the base JSON where users are stored
type UserDatabase struct {
	Users []pwMan.UserInformation `json:"users"`
}

var (
	userDB *UserDatabase
	dbLock = new(sync.RWMutex)
)

// CreateDB creates an empty database
func CreateDB() {
	logrus.Info("Creating empty database")
	dbpath := viper.GetString("databasepath")
	if _, err := os.Stat(dbpath); err == nil {
		// path exists
		logrus.Fatalf("Database %s seems to be already present. Remove it if you really want to start from scratch", dbpath)
	}

	jsonFile, err := os.Create(dbpath)
	if err != nil {
		logrus.Fatal(err)
	}

	jsonFile.Write([]byte("{}"))
	jsonFile.Close()
}

// InitDB From article by Karl Seguin (http://openmymind.net/Golang-Hot-Configuration-Reload/)
func InitDB() {
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
	dbpath := viper.GetString("databasepath")
	temp := new(UserDatabase)

	// @TODO: Set the database path a parameter
	file, e := ioutil.ReadFile(dbpath)
	if e != nil {
		logrus.Warn("Database not found, storing in a new one: ", dbpath)
	} else {
		if err := json.Unmarshal(file, temp); err != nil {
			logrus.Fatal("Error parsing JSON database: ", err)
		}
	}

	dbLock.Lock()
	userDB = temp
	dbLock.Unlock()
}

// AddUser adds a user
func (*UserDatabase) AddUser(u pwMan.UserInformation) string {
	logrus.Info("Adding user to the database")
	db := GetDB()

	// @TODO: Check that it doesn't already exist
	for _, item := range db.Users {
		if item.Username == u.Username {
			logrus.Errorf("User %s already exists in the database", u.Username)
			return "User already present"
		}
	}

	db.Users = append(db.Users, u)
	db.Save()
	return ""
}

// Save saves the JSON db :-)
func (udb *UserDatabase) Save() {
	dbpath := viper.GetString("databasepath")

	jsonFile, err := os.Create(dbpath)
	if err != nil {
		logrus.Fatal(err)
	}

	jsondata, err := json.MarshalIndent(udb, "", "  ")
	if err != nil {
		logrus.Fatal("Error saving user: ", err)
	}

	jsonFile.Write(jsondata)
	jsonFile.Close()
}
