// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/treeder/operator/commands"
	"golang.org/x/net/context"
)

// instancesCmd represents the instances command
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List all instances of an app",
	Long:  `List all instance sof an app, be sure to pass in --name flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		logrus.Infoln("args: ", args)
		// logrus.Infoln("", cmd.Flags)

		if name == "" {
			logrus.Errorln("Name flag required")
			return
		}

		config, err := loadConfig()
		if err != nil {
			return
		}

		commands.Instances(ctx, config, name)
	},
}

func init() {
	RootCmd.AddCommand(instancesCmd)
	instancesCmd.PersistentFlags().StringVar(&name, "name", "", "Name of app.")

}
