package commands

import (
	"github.com/mitchellh/goamz/ec2"
	"github.com/treeder/operator/aws"
	"github.com/treeder/operator/common"
	"golang.org/x/net/context"
)

func GetInstances(ctx context.Context, appname string) ([]*ec2.Instance, error) {
	l := common.Logger(ctx)
	// look up instances, find all instances with tag X
	tags := map[string]string{
		"shortname": appname,
		"tool":      "operator",
	}
	instances, err := aws.GetInstances(tags)
	if err != nil {
		l.WithError(err).Errorln("Error getting instances.")
		return nil, err
	}
	return instances, nil
}
