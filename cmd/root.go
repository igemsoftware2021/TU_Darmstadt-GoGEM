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

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var config Config
var username string
var year int
var timeout int
var teamname string
var wpurl string
var password string
var offset string
var force bool
var clean bool
var insecure bool
var redirect bool

type Config struct {
	URLS            map[string]string `mapstructure:"urls"`
	URLORDER        []string          `mapstructure:"order"`
	FONTS           map[string]string `mapstructure:"fonts"`
	CUSTOMREDIRECTS map[string]string `mapstructure:"customredirects"`
	LOGINURL        string            `mapstructure:"loginurl"`
	LOGOUTURL       string            `mapstructure:"logouturl"`
	PREFIXPAGEURL   string            `mapstructure:"prefixurl"`
	MATHJAXURL      string            `mapstructure:"mathjaxurl"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GoGEM",
	Short: "Upload your Wiki to iGEM",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./GoGEM.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".GoGEM" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("GoGEM")
		viper.SetConfigType("json")

		viper.AddConfigPath(".") // adding current directory as second search path
		viper.SetConfigName("GoGEM")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found, please create GoGEM.json in the current working directory, or the current users home directory.")
		os.Exit(1)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}
}
