package commands

import (
	"fmt"

	"github.com/treeder/operator/aws"
)

type Config struct {
	Aws     *aws.AwsConfig
	Docker  *DockerConfig
	Logging *LoggingConfig
}

type DockerConfig struct {
	Username string
	Password string
}

func (c *DockerConfig) Validate() error {
	if c.Username == "" || c.Password == "" {
		return fmt.Errorf("Missing required docker config")
	}
	return nil
}

type LoggingConfig struct {
	SyslogURL string
}

func (c *LoggingConfig) Validate() error {
	// if c.Username == "" || c.Password == "" {
	// return fmt.Errorf("Missing required docker config")
	// }
	return nil
}
