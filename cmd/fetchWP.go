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

	"github.com/spf13/cobra"
	"goGEM-GoStatic"
)

var project_dir string

// fetchWPCmd represents the fetchWP command
var fetchWPCmd = &cobra.Command{
	Use:   "fetchWP",
	Short: "Clone a WordPress Site to your PC, maintaining all static functionality",
	Long: `Clone a WordPress Site to your PC, maintaining all static functionality.
		Use: goGEM fetchWP [URL]`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloning WordPress Site")
		fmt.Println("URL:", args[0])

		gogemgostatic.GoStatic(args[0], project_dir)
	},
}

func init() {
	rootCmd.AddCommand(fetchWPCmd)

	fetchWPCmd.Flags().StringVarP(&project_dir, "dir", "d", "", "Project Directory")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchWPCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchWPCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
