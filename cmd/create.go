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
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runCreateCmd,
	Args: func(cmd *cobra.Command, args []string) error {
		if file == "" && len(args) < 1 {
			return errors.New("You must supply the type of resource you are creating")
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&file, "file", "", "The file to parse to create stuff")
	rootCmd.AddCommand(createCmd)
}

func runCreateCmd(cmd *cobra.Command, args []string) {
	r, err := getInputStream(args)
	if err != nil {
		log.Fatal(err)
	}

	if err = createResourceReader(r); err != nil {
		log.Fatal(err)
	}
}

func createResourceReader(r io.Reader) error {
	resources, err := loadResources(r)
	if err != nil {
		return err
	}

	for x, l := 0, len(resources); x < l; x++ {
		if err = createResource(resources[x]); err != nil {
			return err
		}

		fmt.Printf("Created %s\n", resources[x].Type)
	}

	return nil
}

func getInputStream(args []string) (io.Reader, error) {
	if file != "" {
		return os.Open(file)
	}

	var additionalArgs []string
	if len(args) > 1 {
		additionalArgs = args[1:]
	}

	rb, err := emptyResourceYaml(args[0], additionalArgs)
	if err != nil {
		return nil, err
	}

	if additionalArgs == nil {
		return editReader(rb)
	}

	return bytes.NewReader(rb), nil
}
