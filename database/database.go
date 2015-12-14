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

package database

import (
	"github.com/lstep/2fanginx/pwMan"
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
	logrus.Info("(init) Loading database of users from ", viper.GetString("databasepath"))
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

// FindUser returns user object if found
func (udb *UserDatabase) FindUser(username string) *pwMan.UserInformation {
	//db := database.GetDB()

	if len(username) > 0 {
		for _, item := range udb.Users {
			if item.Username == username {
				return &item
			}
		}
	}
	return nil
}

// AddUser adds a user
func (udb *UserDatabase) AddUser(u pwMan.UserInformation) string {
	logrus.Infof("Adding user %s to the database", u.Username)

	// @TODO: Check that it doesn't already exist
	for _, item := range udb.Users {
		if item.Username == u.Username {
			logrus.Errorf("User %s already exists in the database", u.Username)
			return "User already present"
		}
	}

	udb.Users = append(udb.Users, u)
	udb.Save()
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
