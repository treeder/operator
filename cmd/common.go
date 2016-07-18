package cmd

import (
	"strings"

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
	// logrus.Infof("awsconfig: %+v ", awsConfig)
	awsConfig.PrivateKey = strings.TrimPrefix(awsConfig.PrivateKey, "\"")
	awsConfig.PrivateKey = strings.TrimSuffix(awsConfig.PrivateKey, "\"")
	// logrus.Infof("awsconfig: %+v ", awsConfig)
	awsConfig.PrivateKey = strings.Replace(awsConfig.PrivateKey, "\\n", "\n", -1)
	// logrus.Infof("awsconfig: %+v ", awsConfig)
	// awsConfig.PrivateKey = strings.Replace(awsConfig.PrivateKey, "-----END RSA PRIVATE KEY-----\n", "", -1)
	// logrus.Infof("awsconfig: %+v ", awsConfig)

	err := awsConfig.Validate()
	if err != nil {
		logrus.WithError(err).Errorln("Invalid environment variables")
		return nil, err
	}

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
