package commands

import (
	"github.com/treeder/operator/aws"
	"github.com/treeder/operator/common"
	"golang.org/x/net/context"
)

func Shell(ctx context.Context, config *Config, name string, sshCmd string) error {
	ctx, l := common.LoggerWithFields(ctx, map[string]interface{}{
		"command": "shell",
		"name":    name,
	})

	instances, err := GetInstances(ctx, name)
	if err != nil {
		l.WithError(err).Errorln("Error getting instance info.")
		return err
	}
	if len(instances) == 0 {
		l.Infoln("No instances running for", name)
		return nil
	}
	for i, instance := range instances {
		ctx, l = common.LoggerWithFields(ctx, map[string]interface{}{"instance_id": instance.InstanceId, "host": instance.DNSName})
		output, err := aws.RunCommandOnServerWithOutput(ctx, config.Aws, sshCmd, instance)
		if err != nil {
			l.WithError(err).Errorln("Error running command")
			return err
		}
		l.Infoln(i, " ssh executed: ", output)
	}
	return nil
}
