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
	"2fanginx/pwMan"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/lstep/gototp"
	"github.com/spf13/cobra"
)

func CreateUser(cmd *cobra.Command, args []string) {

	fmt.Println("Creating User...")

	// Generate TOTP
	init2FA, err := gototp.New(gototp.RandomSecret(10))
	if err != nil {
		logrus.Error(err)
		return
	}
	//fmt.Println(o.QRCodeTerminal("label"))
	logrus.Infof("2FA secret: %s\n", init2FA.Secret())

	pwMan.NewUser("lstep", "foobar", init2FA.Secret())

}
