package awssdk

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	awslocal "github.com/treeder/operator/aws"
	"github.com/treeder/operator/common"
	"github.com/treeder/operator/ssh"
)

func GetEc2() (*ec2.EC2, error) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil, err
	}
	svc := ec2.New(sess)
	return svc, nil
}

// LaunchServer
// TODO: don't take cluster id, should be more generic, maybe just pass in tags?
func LaunchServer(ctx context.Context, config *awslocal.AwsConfig, instanceType string, tags map[string]string) (*ec2.Instance, error) {
	l := common.Logger(ctx)
	e, err := GetEc2()
	if err != nil {
		return nil, err
	}
	imageId := "ami-812ec0ec" // Rancher
	// ImageId:      "ami-8d6485e0", // CoreOS
	ec2Options := &ec2.RunInstancesInput{
		ImageId:      aws.String(imageId),
		MaxCount:     aws.Int64(1), // Required
		MinCount:     aws.Int64(1), // Required
		InstanceType: aws.String(instanceType),
		SubnetId:     aws.String(config.SubnetId),
		KeyName:      aws.String(config.KeyPair),
		// UserData:     userData,
		SecurityGroupIds: []*string{
			aws.String(config.SecurityGroup),
		},
	}
	resp, err := e.RunInstances(ec2Options)
	if err != nil {
		logrus.Errorln("Error running instances", "error", err, "options:", ec2Options)
		return nil, err
	}
	for _, inst := range resp.Instances {
		logrus.Info("Now running", "instance", inst.InstanceId)
	}
	inst := resp.Instances[0]
	l = l.WithField("instance_id", inst.InstanceId)
	logrus.Info("Make sure you terminate instances to stop the cash flow!")

	etags := []*ec2.Tag{}
	for k, v := range tags {
		etags = append(etags, &ec2.Tag{Key: aws.String(k), Value: aws.String(v)})
	}
	_, err = e.CreateTags(&ec2.CreateTagsInput{Tags: etags, Resources: []*string{inst.InstanceId}})
	if err != nil {
		l.WithError(err).Errorln("Error creating tags!")
		// TODO: should terminate instance here or something?
		return nil, err
	}

	// Now we'll wait until it fires up and sshttp is available
	var instance *ec2.Instance
	ticker := time.NewTicker(time.Second * 2)
L:
	for {
		select {
		case <-ticker.C:
			instance, err = GetInstanceInfo(e, *inst.InstanceId)
			if err != nil {
				l.WithError(err).Errorln("Error getting instance info")
				return instance, err
			}
			if checkIfUp(ctx, config, instance) {
				break L
			}
		case <-time.After(5 * time.Minute):
			ticker.Stop()
			l.Warnln("Timed out trying to connect.")
			return instance, fmt.Errorf("Timeout trying to start instance")
		}
	}
	l.Infoln("Instance is running and Docker is online.")
	return instance, err
}

func GetInstanceInfo(e *ec2.EC2, instanceId string) (*ec2.Instance, error) {
	iResp, err := e.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{aws.String(instanceId)}})
	if err != nil {
		logrus.WithError(err).Errorln("Couldn't get instance details")
		return nil, err
	}
	logrus.Debug("GetInstanceInfo", "response", iResp)
	if len(iResp.Reservations) == 0 {
		// instance no longer there
		return nil, fmt.Errorf("Instance not found on aws.")
	}
	instance := iResp.Reservations[0].Instances[0]
	return instance, err
}

func checkIfUp(ctx context.Context, config *awslocal.AwsConfig, i *ec2.Instance) bool {
	logrus.Info("Checking instance status", "state", i.State, "id", i.InstanceId)
	if *i.State.Name != "running" {
		return false
	}
	if *i.PublicDnsName == "" { // wait for it to get a public dns entry (takes a bit)
		return false
	}
	// ssh in and see if docker is alive
	output, err := RunCommandOnServerWithOutput(ctx, config, "docker ps", i)
	if err != nil {
		logrus.WithError(err).Errorln("error excuting ssh command to check if server is up")
		return false
	}
	logrus.Println("output from docker ps:", output)

	return true
}

func RunCommandOnServer(ctx context.Context, config *awslocal.AwsConfig, cmd string, instance *ec2.Instance) error {
	return RunCommandsOnServer2(ctx, config, []string{cmd}, instance, os.Stdout)
}

func RunCommandsOnServer(ctx context.Context, config *awslocal.AwsConfig, cmds []string, instance *ec2.Instance) error {
	return RunCommandsOnServer2(ctx, config, cmds, instance, os.Stdout)
}

func RunCommandsOnServer2(ctx context.Context, config *awslocal.AwsConfig, cmds []string, instance *ec2.Instance, w io.Writer) error {
	log := common.Logger(ctx)
	s, err := opssh.NewSession(*instance.PublicDnsName, config.PrivateKey)
	if err != nil {
		log.WithError(err).Errorln("could not create ssh session!")
		return err
	}
	defer s.Close()

	for _, cmd := range cmds {
		log.Println("run cmd: " + cmd)
		err = s.Run(cmd, w)
		if err != nil {
			log.WithError(err).Errorln("Ssh command failed!", cmd)
			return fmt.Errorf("could not execute command on server: %v", err)
		}
	}
	return nil
}

func RunCommandOnServerWithOutput(ctx context.Context, config *awslocal.AwsConfig, cmd string, instance *ec2.Instance) (string, error) {
	b := &bytes.Buffer{}
	w := io.MultiWriter(b, os.Stdout)
	err := RunCommandsOnServer2(ctx, config, []string{cmd}, instance, w)
	// log.Println("BBBB:", b.String())
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func GetRunningInstances(tags map[string]string) ([]*ec2.Instance, error) {
	filter := &ec2.Filter{Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}}
	return GetInstances(tags, []*ec2.Filter{filter})
}

func GetInstances(tags map[string]string, filters []*ec2.Filter) ([]*ec2.Instance, error) {
	e, err := GetEc2()
	if err != nil {
		return nil, err
	}
	if filters == nil {
		filters = []*ec2.Filter{}
	}
	for k, v := range tags {
		filters = append(filters, &ec2.Filter{Name: aws.String("tag:" + k), Values: []*string{aws.String(v)}})
	}
	resp, err := e.DescribeInstances(&ec2.DescribeInstancesInput{Filters: filters})
	if err != nil {
		return nil, err
	}
	instances := []*ec2.Instance{}
	for _, r := range resp.Reservations {
		for _, inst := range r.Instances {
			instances = append(instances, inst)
		}
	}
	return instances, err
}
