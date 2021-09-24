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
	"bytes"
	"fmt"
	"syscall"

	h "github.com/Jackd4w/goGEM/pkg/Handler"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "CAUTION!! DESTRUCTIVE ACTION! Purge your Wiki Pages from the Server",
	Long: `CAUTION!! DESTRUCTIVE ACTION! Purge your Wiki from the iGEM Servers.
	Files can not be deleted, but pages can be overwritten with no content. Usefull for cleaning up before setting up your actual Wiki.
	THIS IS A DESTRUCTIVE ACTION, you will be required to re enter your password.
	Usage: gogem purge -u "[Username]" -y [Wiki Year] -t "[Teamname]" -o "[Offset]"`,
	Run: func(cmd *cobra.Command, args []string) {
		println("This is a DESTRUCTIVE ACTION, you will be required to re enter your password after you logged in. If you want to abort please hit 'Ctrl + C' on your keyboard or close the shell")
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

		println(fmt.Sprintf("Getting all Pages with prefix %s/%s from https://%d.igem.org", teamname, offset, year))
		pages, err := session.GetAllPages()
		if err != nil {
			println(err.Error())
			return
		}
		for _, page := range pages {
			println(fmt.Sprintf("https://%d.igem.org%s", year, page))
			// session.DeletePage(page)
		}
		println("")
		println("-------------------------------------------------------------")
		println("ARE YOU SURE YOU WANT TO DELETE ALL PAGES ABOVE?")
		println("THIS ACTION CAN NOT BE UNDONE!")
		println("-------------------------------------------------------------")
		print("Re-Enter your password to continue:")
		reEnteredPassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			println(err.Error())
			return
		}
		println("")
		if !bytes.Equal(reEnteredPassword, bytePassword) {
			println("Password Mismatch, aborting...")
			return
		}
		println("")
		println("Purging...")
		for _, page := range pages {
			println(page)
			session.DeletePage(page)
		}
		println("")
		println("Purge complete, logging out")

	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)

	purgeCmd.Flags().StringVarP(&username, "username", "u", "", "Username(required)")
	purgeCmd.MarkFlagRequired("username")
	purgeCmd.Flags().IntVarP(&year, "year", "y", 2021, "Year(required)")
	purgeCmd.MarkFlagRequired("year")
	purgeCmd.Flags().StringVarP(&teamname, "teamname", "t", "", "Teamname(required)")
	purgeCmd.MarkFlagRequired("teamname")

	purgeCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	purgeCmd.Flags().StringVarP(&offset, "offset", "o", "", "Offset from your Teams Namespace root")
	purgeCmd.Flags().StringVarP(&loginURL, "login", "L", "https://igem.org/Login2", "LoginURL, set by default")
	purgeCmd.Flags().StringVarP(&logoutURL, "logout", "l", "https://igem.org/Logout", "LogoutURL, set by default")
	purgeCmd.Flags().StringVarP(&PrefixPageURL, "prefix", "P", fmt.Sprintf("https://%d.igem.org/wiki/index.php?title=Special:PrefixIndex", year), "Special Page 'All Pages with prefix', set by default")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// purgeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// purgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
