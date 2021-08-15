/*
Copyright © 2021 Kai Kabuth <kai.kabuth@stud.tu-darmstadt.de>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	fh "goGEM-FileHandling"
	wp "goGEM-GoStatic"
	h "goGEM-Handler"
)

var username string
var year int
var teamname string
var wpurl string
var password string
var loginURL string
var logoutURL string
var offset string
var force bool
var clean bool

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload your WordPress Page to iGEM",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get necessary data

		println(fmt.Sprintf("Upload %s for %s to https://%d.igem.org/Team:%s", wpurl, username, year, teamname))
		fmt.Print("Enter Password: ")
    	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    	if err != nil {
			println(err.Error())
			return
		}
		println("")
		password = string(bytePassword)

		// Establish connection with iGEM Servers
		println("Logging in...")
		session, err := h.NewHandler(year, username, password, teamname, offset, loginURL, logoutURL)
		if err != nil {
			if err.Error() == "loginFailed" {
				println("Login failed, please try again")
				return
			}
			println(err.Error())
			return
		}
		defer session.Logout()
		println("Logged in")
		// Clone WordPress Page
		println("Cloning WordPress Page...")
		project_path, err := wp.GoStatic(wpurl, "")
		if err != nil {
			println(err.Error())
			return
		}
		defer cleanUp(project_path)
		println("Cloning successfull, begining upload...")
		// Prepare Files and Upload them

		if err := fh.PrepFilesForIGEM(teamname, project_path, session); err != nil {
			println(err.Error())
			return
		}
		println("Upload successfull")
		println("Logging out")

	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(&username, "username", "u", "", "Username(required)")
	uploadCmd.MarkFlagRequired("username")
	uploadCmd.Flags().IntVarP(&year, "year", "y", 2021, "Year(required)")
	uploadCmd.MarkFlagRequired("year")
	uploadCmd.Flags().StringVarP(&teamname, "teamname", "t", "", "Teamname(required)")
	uploadCmd.MarkFlagRequired("teamname")
	uploadCmd.Flags().StringVarP(&wpurl, "wpurl", "w", "", "WordPress URL(required)")
	uploadCmd.MarkFlagRequired("wpurl")
	uploadCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	uploadCmd.Flags().StringVarP(&offset, "offset", "o", "", "Offset from your Teams Namespace root")
	uploadCmd.Flags().StringVarP(&loginURL, "login", "L", "https://igem.org/Login2", "LoginURL, set by default")
	uploadCmd.Flags().StringVarP(&logoutURL, "logout", "l", "https://igem.org/Logout", "LogoutURL, set by default")
	uploadCmd.Flags().BoolVarP(&force, "force", "f", false, "Force")
	uploadCmd.Flags().BoolVarP(&clean, "clean", "c", true, "Clean")
	// uploadCmd.Flags().StringVarP(&password, "password", "p", "", "Password(required)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cleanUp(project_dir string) {
	if clean {
		os.RemoveAll(project_dir)
	}
}