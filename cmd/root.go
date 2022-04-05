/*
Copyright © 2022 Arizona Hanson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/arizonahanson/oryx/pkg/eval"
	"github.com/arizonahanson/oryx/pkg/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "oryx",
	Short:   "Embedded expression language",
	Long:    `Embedded expression language.`,
	Args:    cobra.ExactArgs(1),
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		env := lib.BaseEnv(nil)
		val, err := eval.EvalBytes([]byte(strings.Join(args, " ")), env)
		cobra.CheckErr(err)
		fmt.Println(val)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default ~/.oryx.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")
	if cfgFile == "" {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		// Search config in home directory with filename ".oryx.toml".
		viper.AddConfigPath(home)
		viper.SetConfigName(".oryx")
		// create default config
		viper.SafeWriteConfig()
	} else {
		// Use specified file
		viper.SetConfigFile(cfgFile)
		// create if not exists
		viper.SafeWriteConfigAs(cfgFile)
	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// found but could not read?
			cobra.CheckErr(err)
		}
	}
	cfgFile = viper.ConfigFileUsed()
	// watch the config file for changes
	viper.WatchConfig()
	// read in environment variables that match.
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("ORYX")
	// cast types on `Get` to match default values
	viper.SetTypeByDefaultValue(true)
}

func Quit() {
	os.Exit(0)
}
