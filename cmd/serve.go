// Copyright © 2015 NAME HERE <EMAIL ADDRESS>
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
	"2fanginx/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd respresents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch the server",
	//Long: ``,
	Run: server.Run,
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// @NOTE: Do not set the default values here, that doesn't work correctly!
	serveCmd.Flags().BoolP("daemon", "d", false, "Run as a daemon and detach from terminal")
	serveCmd.Flags().StringP("address", "a", "", "Set address:port to listen on")
	serveCmd.Flags().StringP("domain", "e", "", "Set domain to use with cookies")
	serveCmd.Flags().IntP("cookiemaxage", "m", 4, "Set magical cookie's lifetime before it expires (in hours)")
	serveCmd.Flags().StringP("cookiesecret", "c", "", "Secret string to use in signing cookies")

	viper.BindPFlag("address", serveCmd.Flags().Lookup("address"))
	viper.BindPFlag("domain", serveCmd.Flags().Lookup("domain"))
	viper.BindPFlag("cookiemaxage", serveCmd.Flags().Lookup("cookiemaxage"))
	viper.BindPFlag("cookiesecret", serveCmd.Flags().Lookup("cookiesecret"))

	viper.SetDefault("address", "127.0.0.1:9434")
	viper.SetDefault("domain", ".secure.mydomain.eu")
	viper.SetDefault("cookiemaxage", 4)
	viper.SetDefault("cookiesecret", "CHOOSE-A-SECRET-YOURSELF")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
