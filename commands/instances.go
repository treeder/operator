package commands

import (
	"github.com/treeder/operator/common"
	"golang.org/x/net/context"
)

func Instances(ctx context.Context, config *Config, name string) error {
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
		ctx, l = common.LoggerWithFields(ctx, map[string]interface{}{
			"instance_id": instance.InstanceId,
			"host":        instance.DNSName,
			"launch_time": instance.LaunchTime,
			"type":        instance.InstanceType,
		})
		l.Infoln(i, " instance ")
	}
	return nil
}
