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

package cmd

import (
	"github.com/lstep/2fanginx/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd respresents the serve command
var createUserCmd = &cobra.Command{
	Use:   "createuser",
	Short: "Create a user",
	//Long: ``,
	Run: server.CreateUser,
}

func init() {
	RootCmd.AddCommand(createUserCmd)

	// @NOTE: Do not set the default values here, that doesn't work correctly!
	createUserCmd.Flags().StringP("name", "n", "", "username")
	createUserCmd.Flags().StringP("hmackey", "p", "", "hmackey")

	viper.BindPFlag("name", createUserCmd.Flags().Lookup("name"))
	viper.BindPFlag("hmackey", createUserCmd.Flags().Lookup("hmackey"))

	viper.SetDefault("hmackey", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
}
