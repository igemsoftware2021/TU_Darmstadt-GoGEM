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
	"log"

	cc "github.com/Jackd4w/goGEM/pkg/checkCriteria"
	"github.com/spf13/cobra"
)

var url bool

var checkCriteriaCMD = &cobra.Command{
	Use:   "checkCriteria",
	Short: "Check Medal Criteria URLs",
	Long:  `Check if your Wiki fullfils the URL Criteria for Medals.
			Award and MedalCriteria are defined in ...
			Usage: gogem checkcriteria -y [year] -t [teamname] -u [boolean]`,

	Run: func(cmd *cobra.Command, args []string) {
		results, err := cc.CheckCriteria(config.URLORDER, config.URLS, teamname, year, url)
		if err != nil {
			log.Fatal(err)
		}

		println("This is only a guideline! URLs and Criteria may not be up to date!")
		println("ABSOLUTELY NO WARRENTY FOR THE GENERATED RESULTS")

		println(results)
	},
}

func init() {
	rootCmd.AddCommand(checkCriteriaCMD)

	checkCriteriaCMD.Flags().IntVarP(&year, "year", "y", 2021, "Year(required)")
	checkCriteriaCMD.MarkFlagRequired("year")
	checkCriteriaCMD.Flags().StringVarP(&teamname, "teamname", "t", "", "Teamname(required)")
	checkCriteriaCMD.MarkFlagRequired("teamname")
	checkCriteriaCMD.Flags().BoolVarP(&url, "url", "u", false, "Print URLs")

}
