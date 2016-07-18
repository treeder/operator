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
	"strings"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/treeder/operator/commands"
)

var name string
var envVars []string
var deployOptions *commands.DeployOptions

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		logrus.Infoln("args: ", args)
		logrus.Infoln("env vars: ", envVars, "len: ", len(envVars))
		// logrus.Infoln("", cmd.Flags)

		if name == "" {
			logrus.Errorln("Name flag required")
			return
		}

		config, err := loadConfig()
		if err != nil {
			return
		}

		image := args[0]

		// The env vars aren't being parsed properly, they're coming in concatenated in the first element of the slice, so this will handle both hopefully
		envVarsMap := map[string]string{}
		for _, v := range envVars {
			logrus.Infoln(v)
			sp := strings.SplitN(v, "=", 2)
			logrus.Infoln(sp)
			envVarsMap[sp[0]] = sp[1]
		}
		deployOptions.EnvVars = envVarsMap
		commands.Deploy(ctx, config, name, image, deployOptions)
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deployCmd.PersistentFlags().StringVar(&name, "name", "", "Name of app.")
	deployCmd.PersistentFlags().StringSliceVarP(&envVars, "env", "e", []string{}, "Environment variables.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployOptions = &commands.DeployOptions{}
	deployCmd.Flags().BoolVar(&deployOptions.Privileged, "privileged", false, "Run in privileged mode.")

}
