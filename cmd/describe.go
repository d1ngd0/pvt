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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)
		}

		t := normalizeType(args[0])
		res, err := getResource(viper.GetInt("project"), t, id)
		if err != nil {
			log.Fatal(err)
		}

		et, err := toType(t)
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(res.Spec, et); err != nil {
			log.Fatal(err)
		}

		if err = render(fmt.Sprintf("describe-%s", t), os.Stdout, et); err != nil {
			log.Fatal(err)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("you must supply a type and an id")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
