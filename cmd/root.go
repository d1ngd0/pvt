/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/d1ngd0/pvt/pivotal"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	c       *pivotal.Client
	file    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pvt",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: myworkCmd.Run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pvt.yaml)")

	rootCmd.PersistentFlags().Int("project", 0, "The project to to pull information from")
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))

	rootCmd.PersistentFlags().String("me", "", "Your username, or someone elses if you want to see their work")
	viper.BindPFlag("me", rootCmd.PersistentFlags().Lookup("me"))

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
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pvt" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(fmt.Sprintf("%s/.config/", home))
		viper.AddConfigPath("/etc/pvt/")

		viper.SetConfigName("pvt")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
	c = pivotal.New(viper.GetString("token"))
}

func editReader(initb []byte) (io.Reader, error) {
	tmpFile, err := ioutil.TempFile("", "pvt-*.yml")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	initMd5 := md5.Sum(initb)
	_, err = tmpFile.Write(initb)
	tmpFile.Close()

	if err != nil {
		return nil, err
	}

	if err = editor(tmpFile.Name()); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(tmpFile.Name())
	if initMd5 == md5.Sum(b) {
		return nil, errors.New("no changes made")
	}
	return bytes.NewReader(b), err
}

func editor(filepath string) error {
	cmd := exec.Command(getEditor(), filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	return editor
}
