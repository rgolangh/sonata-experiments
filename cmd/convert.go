/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"os"
	"strings"

	"github.com/rgolangh/sonata-experiments/internal/backstage"
	"github.com/rgolangh/sonata-experiments/internal/sonata"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [filename]",
	Short: "convert from backstage template to sonataflow definition",
	Long:  `Convert a backstage softare template to an orchestrator sonataflow json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Supply a file name to parse")
		}

		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("Error:", err)
		}
		defer file.Close()

		// Decode YAML into Template struct
		var t backstage.Template
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&t); err != nil {
			return fmt.Errorf("Error:", err)
		}
		flow := sonata.NewFrom(t)

		j, err := json.Marshal(&flow)
		if err != nil {
			return fmt.Errorf("error when marshaling json %v\n", err)
		}

		// convert from json to yaml .The reason I dont marshal directly
		// to yaml is because the struct tags on sonata types are json so
		// all the omitempty are not respected.
		jdoc := make(map[string]interface{})
		json.Unmarshal(j, &jdoc)

		// add id and name under the root json, as those are missing from the sonata operator model
		jdoc["id"] = strings.ReplaceAll(strings.ToLower(t.Metadata.Name), " ", "-")
		jdoc["name"] = t.Metadata.Name
		jdoc["description"] = t.Metadata.Description

		marshal, err := yaml.Marshal(&jdoc)
		if err != nil {
			return fmt.Errorf("failed marshaling yaml%v", err)
		}
		fmt.Println(string(marshal))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
