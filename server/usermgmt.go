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
	"2fanginx/database"
	"2fanginx/pwMan"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/Songmu/prompter"
	"github.com/lstep/gototp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CreateUser is a procedure for creating a user
func CreateUser(cmd *cobra.Command, args []string) {
	database.InitDB()
	db := database.GetDB()

	username := viper.GetString("name")
	if username == "" {
		fmt.Println("Required 'name' parameter not specified")
		return
	}

	fmt.Printf("Creating User %s...\n", username)

	for _, item := range db.Users {
		if item.Username == username {
			logrus.Errorf("User %s already exists in the database", username)
			return
		}
	}

	// Generate TOTP
	init2FA, err := gototp.New(gototp.RandomSecret(10))
	if err != nil {
		logrus.Error(err)
		return
	}

	// check if password was specified, otherwise, go interactive
	password := viper.GetString("password")
	if password == "" {
		password = prompter.Password("Enter password to use")
	}

	user, err := pwMan.NewUser(username, password, init2FA.Secret())
	if err != nil {
		fmt.Printf("Error while creating user %s: %v\n", username, err)
	}

	db.AddUser(*user)

	fmt.Printf("User %s created. Caracteristics :\n", username)
	fmt.Printf("2FA init: %s || QRCode link: %s\n", init2FA.Secret(), init2FA.QRCodeGoogleChartsUrl("Code", 320))
	//fmt.Println(.QRCodeTerminal("label"))

}
