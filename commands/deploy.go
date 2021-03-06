package commands

import (
	"bytes"
	"fmt"

	ec22 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/treeder/operator/aws/awssdk"
	"github.com/treeder/operator/common"
	"golang.org/x/net/context"
)

const defaultInstanceType = "t2.small"

type DeployOptions struct {
	Privileged bool
	EnvVars    map[string]string
}

func Deploy(ctx context.Context, config *Config, name, image string, options *DeployOptions) ([]*ec22.Instance, error) {
	ctx, l := common.LoggerWithFields(ctx, map[string]interface{}{
		"command": "deploy",
		"name":    name,
		"image":   image,
	})
	l.Infof("deployOptions: %+v\n", options)

	// look up instances, find all instances with tag X
	instances, err := GetInstances2(ctx, name)
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
		instance, err := awssdk.LaunchServer(ctx, config.Aws, instanceType, tags)
		if err != nil {
			l.WithError(err).Errorln("Error launching server.")
			return nil, err
		}
		instances = append(instances, instance)
		ctx, l = common.LoggerWithField(ctx, "instance_id", *instance.InstanceId)
		l.Println(instance)

		cmds := []string{}
		l.Println("SYSLOG_URL: ", config.Logging)
		if config.Logging.SyslogURL != "" {
			cmds = append(cmds, logspoutCmd(config.Logging))
		}
		cmds = append(cmds,
			pullCmd(image, config.Docker),
			runCmd(instance, name, image, options))
		err = awssdk.RunCommandsOnServer(ctx, config.Aws, cmds, instance)
		if err != nil {
			l.WithError(err).Errorln("Error starting container")
			return nil, err
		}
		// TOOD: verify the runner container started OK and is running
	} else {
		l.Infoln("Instances already running, updating...")
		for i, instance := range instances {
			ctx, l = common.LoggerWithField(ctx, "instance_id", *instance.InstanceId)
			l.Infoln("Updating instance ", i)
			cmds := []string{}
			// l.Println("SYSLOG_URL: ", config.Logging)
			// if config.Logging.SyslogURL != "" {
			// 	cmds = append(cmds, logspoutCmd(config.Logging))
			// }
			cmds = append(cmds,
				pullCmd(image, config.Docker),
				fmt.Sprintf("docker stop %v", name),
				fmt.Sprintf("docker rm %v", name),
				runCmd(instance, name, image, options))

			// TODO: change this to docker pull, docker stop, then docker run again
			err := awssdk.RunCommandsOnServer(ctx, config.Aws, cmds, instance)
			if err != nil {
				l.WithError(err).Errorln("Error starting iron/runner container")
				return nil, err
			}
		}
	}
	for _, instance := range instances {
		l.Infoln(instance.PublicDnsName)
	}
	return instances, err
}

// this uses itself to use the Docker API on the remote server to pull the image.
func pullCmd(image string, cfg *DockerConfig) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("docker run -v /var/run/docker.sock:/var/run/docker.sock --rm treeder/operator pull %v", image))
	if cfg != nil {
		buffer.WriteString(fmt.Sprintf(" -u %v -p %v", cfg.Username, cfg.Password))
	}
	return buffer.String()
}

func runCmd(instance *ec22.Instance, name, image string, options *DeployOptions) string {
	var buffer bytes.Buffer
	buffer.WriteString("docker run -d ")
	buffer.WriteString(fmt.Sprintf("--name %v ", name))
	if options.Privileged {
		buffer.WriteString("--privileged ")
	}
	// TODO: allow user to set set port, etc
	buffer.WriteString("-p 80:8080 -e PORT=8080 ")
	for k, v := range options.EnvVars {
		buffer.WriteString(fmt.Sprintf("-e \"%v=%v\" ", k, v))
	}
	buffer.WriteString(image)
	return buffer.String()
}

func logspoutCmd(loggingConfig *LoggingConfig) string {
	var buffer bytes.Buffer
	buffer.WriteString("docker run -d --name=logspout --volume=/var/run/docker.sock:/var/run/docker.sock gliderlabs/logspout syslog+udp://")
	buffer.WriteString(loggingConfig.SyslogURL)
	return buffer.String()

}
