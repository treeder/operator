package commands

import (
	ec22 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/goamz/ec2"
	"github.com/treeder/operator/aws"
	"github.com/treeder/operator/aws/awssdk"
	"golang.org/x/net/context"
)

func GetInstances(ctx context.Context, appname string) ([]*ec2.Instance, error) {
	// look up instances, find all instances with tag X
	tags := map[string]string{
		"shortname": appname,
		"tool":      "operator",
	}
	instances, err := aws.GetRunningInstances(tags)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func GetInstances2(ctx context.Context, appname string) ([]*ec22.Instance, error) {
	// look up instances, find all instances with tag X
	tags := map[string]string{
		"shortname": appname,
		"tool":      "operator",
	}
	instances, err := awssdk.GetRunningInstances(tags)
	if err != nil {
		return nil, err
	}
	return instances, nil
}
