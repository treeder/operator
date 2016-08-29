package commands

import (
	"github.com/treeder/operator/aws"
	"github.com/treeder/operator/common"
	"golang.org/x/net/context"
)

func Kill(ctx context.Context, name, instanceId string) error {
	ctx, l := common.LoggerWithFields(ctx, map[string]interface{}{
		"command":     "kill",
		"name":        name,
		"instance_id": instanceId,
	})

	err := aws.KillServer(ctx, instanceId)
	if err != nil {
		l.WithError(err).Errorln("Error running command")
		return err
	}
	l.Infoln(" terminated server ")

	return nil
}
