package commands

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/goamz/ec2"
	"github.com/treeder/operator/aws"
	"github.com/treeder/operator/common"
	"golang.org/x/net/context"
)

const defaultInstanceType = "m4.large"

func Deploy(ctx context.Context, config *Config, name, image string) ([]*ec2.Instance, error) {
	ctx, l := common.LoggerWithFields(ctx, map[string]interface{}{
		"command": "deploy",
		"name":    name,
		"image":   image,
	})

	// look up instances, find all instances with tag X
	instances, err := GetInstances(ctx, name)
	if err != nil {
		l.WithError(err).Errorln("Error getting instance info.")
		return nil, err
	}

	if len(instances) == 0 {
		l.Infoln("No running instances, starting new one.")
		tags := map[string]string{
			"shortname": name,
			"tool":      "operator",
		}
		tags["Name"] = name
		instanceType := defaultInstanceType
		instance, err := aws.LaunchServer(ctx, config.Aws, instanceType, tags)
		if err != nil {
			l.WithError(err).Errorln("Error launching server.")
			return nil, err
		}
		instances = append(instances, instance)
		ctx, l = common.LoggerWithField(ctx, "instance_id", instance.InstanceId)
		l.Println(instance)

		// TODO:		startMonitoringContainers(ctx, instance)
		cmds := []string{
			pullCmd(image, config.Docker),
			runCmd(instance, image),
		}
		err = aws.RunCommandsOnServer(ctx, config.Aws, cmds, instance)
		if err != nil {
			l.WithError(err).Errorln("Error starting container")
			return nil, err
		}
		// TOOD: verify the runner container started OK and is running
	} else {
		l.Infoln("Instances already running, updating...")
		for i, instance := range instances {
			ctx, l = common.LoggerWithField(ctx, "instance_id", instance.InstanceId)
			l.Infoln("Updating instance ", i)
			cmds := []string{
				pullCmd(image, config.Docker),
				fmt.Sprintf("docker stop %v", image),
				runCmd(instance, image),
			}
			// TODO: change this to docker pull, docker stop, then docker run again
			err := aws.RunCommandsOnServer(ctx, config.Aws, cmds, instance)
			if err != nil {
				l.WithError(err).Errorln("Error starting iron/runner container")
				return nil, err
			}
		}
	}
	return instances, err
}

// this uses itself to use the Docker API on the remote server to pull the image.
func pullCmd(image string, cfg *DockerConfig) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("docker run  -v /var/run/docker.sock:/var/run/docker.sock --rm treeder/operator pull %v", image))
	if cfg != nil {
		buffer.WriteString(fmt.Sprintf(" -u %v -p %v", cfg.Username, cfg.Password))
	}
	return buffer.String()
}

func runCmd(instance *ec2.Instance, image string) string {
	var buffer bytes.Buffer
	buffer.WriteString("docker run -d ")
	// TODO: allow user to set set port, etc
	buffer.WriteString("-p 80:8080 -e PORT=8080 ")
	// TODO: add env vars
	buffer.WriteString(image)
	return buffer.String()
}
