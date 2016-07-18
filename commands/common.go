package commands

import (
	"github.com/mitchellh/goamz/ec2"
	"github.com/treeder/operator/aws"
	"golang.org/x/net/context"
)

func GetInstances(ctx context.Context, appname string) ([]*ec2.Instance, error) {
	// look up instances, find all instances with tag X
	tags := map[string]string{
		"shortname": appname,
		"tool":      "operator",
	}
	instances, err := aws.GetInstances(tags)
	if err != nil {
		return nil, err
	}
	return instances, nil
}
