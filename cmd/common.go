package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/treeder/operator/aws"
	"github.com/treeder/operator/commands"
)

func loadConfig() (*commands.Config, error) {
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
		return nil, err
	}
	// logrus.Infof("awsconfig: %+v ", awsConfig)

	dockerConfig := &commands.DockerConfig{
		Username: viper.GetString("DOCKER_USERNAME"),
		Password: viper.GetString("DOCKER_PASSWORD"),
	}
	dockerConfig.Validate()
	if err != nil {
		logrus.WithError(err).Errorln("Invalid environment variables")
		return nil, err
	}
	return &commands.Config{Aws: awsConfig, Docker: dockerConfig}, nil

}
