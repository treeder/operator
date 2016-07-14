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

// var sshCmd string

// shCmd represents the sh command
var shCmd = &cobra.Command{
	Use:   "sh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		logrus.Infoln("args: ", args)
		// logrus.Infoln("", cmd.Flags)

		config, err := loadConfig()
		if err != nil {
			return
		}

		commands.Shell(ctx, config, name, args[0])
	},
}

func init() {
	RootCmd.AddCommand(shCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	shCmd.PersistentFlags().StringVar(&name, "name", "", "Name of app.")
	// shCmd.PersistentFlags().StringVar(&sshCmd, "ssh", "", "runs ssh command on each")

}
