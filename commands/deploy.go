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

func Deploy(ctx context.Context, config *aws.AwsConfig, name, image string) ([]*ec2.Instance, error) {
	ctx, l := common.LoggerWithFields(ctx, map[string]interface{}{
		"command": "deploy",
		"image":   image,
	})

	// look up instances, find all instances with tag X
	tags := map[string]string{
		"shortname": name,
		"tool":      "operator",
	}
	instances, err := aws.GetInstances(tags)
	if err != nil {
		l.WithError(err).Errorln("Error getting instance info.")
		return nil, err
	}

	if len(instances) == 0 {
		l.Infoln("No running instances, starting new one.")
		tags["Name"] = name
		instanceType := defaultInstanceType
		instance, err := aws.LaunchServer(ctx, config, instanceType, tags)
		if err != nil {
			l.WithError(err).Errorln("Error launching server.")
			return nil, err
		}
		instances = append(instances, instance)
		ctx, l = common.LoggerWithField(ctx, "instance_id", instance.InstanceId)
		l.Println(instance)

		// TODO:		startMonitoringContainers(ctx, instance)
		runc := runCmd(instance, image)
		l.Println("Running: ", runc)
		output, err := aws.RunCommandOnServerWithOutput(ctx, config, runc, instance)
		if err != nil {
			l.WithError(err).Errorln("Error starting container")
			return nil, err
		}
		l.Infoln("start iron/runner output: ", output)
		// TOOD: verify the runner container started OK and is running
	} else {
		l.Infoln("Instances already running, updating...")

		for i, instance := range instances {
			ctx, l = common.LoggerWithField(ctx, "instance_id", instance.InstanceId)
			l.Infoln("Updating instance ", i)
			cmds := []string{
				fmt.Sprintf("docker pull %v", image),
				fmt.Sprintf("docker stop %v", image),
				runCmd(instance, image),
			}
			// TODO: change this to docker pull, docker stop, then docker run again
			err := aws.RunCommandsOnServer(ctx, config, cmds, instance)
			if err != nil {
				l.WithError(err).Errorln("Error starting iron/runner container")
				return nil, err
			}
		}
	}
	return instances, err
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
