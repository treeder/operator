package aws

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/ec2"
	"github.com/treeder/operator/common"
	"github.com/treeder/operator/ssh"
	"gopkg.in/inconshreveable/log15.v2"
)

type AwsConfig struct {
	SubnetId      string
	KeyPair       string
	PrivateKey    string
	SecurityGroup string
}

func (c *AwsConfig) Validate() error {
	if c.SubnetId == "" || c.KeyPair == "" || c.PrivateKey == "" || c.SecurityGroup == "" {
		return fmt.Errorf("Missing required aws config")
	}
	return nil
}

func GetEc2() (*ec2.EC2, error) {
	// auth, err := aws.GetAuth(j.Opts.AwsAccessKey, j.Opts.AwsSecretKey)
	auth, err := aws.EnvAuth()
	if err != nil {
		log.WithError(err).Errorln("Error aws.GetAuth")
		return nil, err
	}
	e := ec2.New(auth, aws.USEast)
	return e, nil
}

func GetInstanceInfo(e *ec2.EC2, instanceId string) (*ec2.Instance, error) {
	iResp, err := e.Instances([]string{instanceId}, nil)
	if err != nil {
		log15.Crit("Couldn't get instance details", "error", err)
		return nil, err
	}
	log15.Debug("GetInstanceInfo", "response", iResp)
	if len(iResp.Reservations) == 0 {
		// instance no longer there
		return nil, fmt.Errorf("Instance not found on aws.")
	}
	instance := iResp.Reservations[0].Instances[0]
	return &instance, err
}

func GetInstances(tags map[string]string) ([]*ec2.Instance, error) {
	e, err := GetEc2()
	if err != nil {
		return nil, err
	}
	filter := ec2.NewFilter()
	for k, v := range tags {
		filter.Add("tag:"+k, v)
	}
	resp, err := e.Instances(nil, filter)
	if err != nil {
		return nil, err
	}
	instances := []*ec2.Instance{}
	for _, r := range resp.Reservations {
		for _, inst := range r.Instances {
			instances = append(instances, &inst)
		}
	}
	return instances, err
}

func KillServer(ctx context.Context, instanceId string) error {
	l := common.Logger(ctx)
	l.Infoln("terminating")
	e, err := GetEc2()
	if err != nil {
		return err
	}
	_, err = e.TerminateInstances([]string{instanceId})
	return err
}

// LaunchServer
// TODO: don't take cluster id, should be more generic, maybe just pass in tags?
func LaunchServer(ctx context.Context, config *AwsConfig, instanceType string, tags map[string]string) (*ec2.Instance, error) {
	l := common.Logger(ctx)
	e, err := GetEc2()
	if err != nil {
		return nil, err
	}

	ec2Options := ec2.RunInstances{
		// ImageId:      "ami-8d6485e0", // CoreOS
		ImageId:      "ami-812ec0ec", // "ami-53045239", // Rancher
		InstanceType: instanceType,
		SubnetId:     config.SubnetId,
		KeyName:      config.KeyPair,
		// UserData:     userData,
		SecurityGroups: []ec2.SecurityGroup{
			{Id: config.SecurityGroup},
		},
	}
	resp, err := e.RunInstances(&ec2Options)
	if err != nil {
		log.Errorln("Error running instances", "error", err, "options:", ec2Options)
		return nil, err
	}
	for _, inst := range resp.Instances {
		log15.Info("Now running", "instance", inst.InstanceId)
	}
	inst := resp.Instances[0]
	l = l.WithField("instance_id", inst.InstanceId)
	log15.Info("Make sure you terminate instances to stop the cash flow!")

	etags := []ec2.Tag{}
	for k, v := range tags {
		etags = append(etags, ec2.Tag{Key: k, Value: v})
	}
	_, err = e.CreateTags([]string{inst.InstanceId}, etags)
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
			instance, err = GetInstanceInfo(e, inst.InstanceId)
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

func checkIfUp(ctx context.Context, config *AwsConfig, i *ec2.Instance) bool {
	log.Info("Checking instance status", "state", i.State, "id", i.InstanceId)
	if i.State.Name != "running" {
		return false
	}
	if i.DNSName == "" { // wait for it to get a public dns entry (takes a bit)
		return false
	}
	// ssh in and see if docker is alive
	output, err := RunCommandOnServerWithOutput(ctx, config, "docker ps", i)
	if err != nil {
		log.WithError(err).Errorln("error excuting ssh command to check if server is up")
		return false
	}
	log.Println("output from docker ps:", output)

	return true
}

func RunCommandOnServer(ctx context.Context, config *AwsConfig, cmd string, instance *ec2.Instance) error {
	return RunCommandsOnServer2(ctx, config, []string{cmd}, instance, os.Stdout)
}

func RunCommandsOnServer(ctx context.Context, config *AwsConfig, cmds []string, instance *ec2.Instance) error {
	return RunCommandsOnServer2(ctx, config, cmds, instance, os.Stdout)
}

func RunCommandsOnServer2(ctx context.Context, config *AwsConfig, cmds []string, instance *ec2.Instance, w io.Writer) error {
	l := common.Logger(ctx)
	s, err := opssh.NewSession(instance.DNSName, config.PrivateKey)
	if err != nil {
		log.WithError(err).Errorln("could not create ssh session!")
		return err
	}
	defer s.Close()

	for _, cmd := range cmds {
		l.Println("run cmd: " + cmd)
		err = s.Run(cmd, w)
		if err != nil {
			l.WithError(err).Errorln("Ssh command failed!", cmd)
			return fmt.Errorf("could not execute command on server: %v", err)
		}
	}
	return nil
}

func RunCommandOnServerWithOutput(ctx context.Context, config *AwsConfig, cmd string, instance *ec2.Instance) (string, error) {
	b := &bytes.Buffer{}
	w := io.MultiWriter(b, os.Stdout)
	err := RunCommandsOnServer2(ctx, config, []string{cmd}, instance, w)
	// log.Println("BBBB:", b.String())
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
