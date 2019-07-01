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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runUpdateCmd,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("You must supply the type of resource you are updating and an ID")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdateCmd(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal(err)
	}

	res, err := getResource(viper.GetInt("project"), args[0], id)
	if err != nil {
		log.Fatal(err)
	}

	r, err := getUpdateInputStream(res, args)
	if err != nil {
		log.Fatal(err)
	}

	if err = updateResourceReader(r); err != nil {
		log.Fatal(err)
	}
}

func updateResourceReader(r io.Reader) error {
	resources, err := loadResources(r)
	if err != nil {
		return err
	}

	for x, l := 0, len(resources); x < l; x++ {
		if err = updateResource(resources[x]); err != nil {
			return err
		}

		fmt.Printf("Updated %s\n", resources[x].Type)
	}

	return nil
}

func getUpdateInputStream(res resource, args []string) (io.Reader, error) {
	et, err := toType(res.Type)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(res.Spec, et); err != nil {
		return nil, err
	}

	if len(args) > 2 {
		applyArguments(et, args[2:])
	}

	if res.Spec, err = json.Marshal(et); err != nil {
		return nil, err
	}

	rb, err := resourceYaml(res)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) <= 2 {
		return editReader(rb)
	}

	return bytes.NewReader(rb), nil
}
