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
	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/treeder/operator/aws"
	"github.com/treeder/operator/commands"
)

var name string

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
		// logrus.Infoln("", cmd.Flags)

		// Get required env vars:
		awsConfig := &aws.AwsConfig{}
		awsConfig.SubnetId = viper.GetString("AWS_SUBNET_ID")
		awsConfig.KeyPair = viper.GetString("AWS_KEY_PAIR")
		// UserData:     userData,
		awsConfig.SecurityGroup = viper.GetString("AWS_SECURITY_GROUP")
		awsConfig.PrivateKey = viper.GetString("AWS_PRIVATE_KEY")
		err := awsConfig.Validate()
		if err != nil {
			logrus.WithError(err).Errorln("Invalid environment variables")
			return
		}
		logrus.Infof("awsconfig: %+v ", awsConfig)

		dockerConfig := &commands.DockerConfig{
			Username: viper.GetString("DOCKER_USERNAME"),
			Password: viper.GetString("DOCKER_PASSWORD"),
		}
		dockerConfig.Validate()
		if err != nil {
			logrus.WithError(err).Errorln("Invalid environment variables")
			return
		}

		commands.Deploy(ctx, &commands.Config{Aws: awsConfig, Docker: dockerConfig}, name, args[0])
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deployCmd.PersistentFlags().StringVar(&name, "name", "", "Name of app.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
