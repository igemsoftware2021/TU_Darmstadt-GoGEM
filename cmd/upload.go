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

	fh "github.com/Jackd4w/goGEM/pkg/FileHandling"
	wp "github.com/Jackd4w/goGEM/pkg/GoStatic"
	h "github.com/Jackd4w/goGEM/pkg/Handler"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload your WordPress Page to iGEM",
	Long: `Curls every URL that is reachable from the specified entry URL.
	Replaces every relative Link on the WP-Page with static links pointing to the iGEM Servers.
	If you want to clone your Wiki to https://2021.igem.org/Team:TU_Darmstadt/test/[...] then the command would be:
	gogem upload -u "[Your Username]" -y 2021 -t "TU_Darmstadt" -w "[Your WP Wiki]" -o "test"` ,

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
		session, err := h.NewHandler(year, username, password, teamname, offset, loginURL, logoutURL, PrefixPageURL)
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
		project_path, err := wp.GoStatic(wpurl, "", insecure)
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
	uploadCmd.Flags().StringVarP(&PrefixPageURL, "prefix", "P", fmt.Sprintf("https://%d.igem.org/wiki/index.php?title=Special:PrefixIndex", year), "Special Page 'All Pages with prefix', set by default")
	uploadCmd.Flags().BoolVarP(&force, "force", "f", false, "Force")
	uploadCmd.Flags().BoolVarP(&clean, "clean", "c", true, "Clean")
	uploadCmd.Flags().BoolVarP(&insecure, "insecure", "i", false, "Ignores HTTPS Certificate warnings")
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