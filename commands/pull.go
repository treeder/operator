package commands

import (
	"encoding/base64"
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/treeder/operator/common"
	"golang.org/x/net/context"
)

func Pull(ctx context.Context, image, username, password string) error {
	ctx, log := common.LoggerWithFields(ctx, map[string]interface{}{
		"command": "pull",
	})

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.WithError(err).Fatalln("Error making docker client!")
		return err
	}

	repo, tag := docker.ParseRepositoryTag(image)
	log.Infoln("username:", username, "password:", password, "image:", image, "repo:", repo, "tag:", tag)

	auth := ""
	if username != "" {
		auth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	}

	authConfig := docker.AuthConfiguration{}
	if auth != "" {
		read := strings.NewReader(fmt.Sprintf(`{"docker.io":{"auth":"%s"}}`, auth))
		ac, err := docker.NewAuthConfigurations(read)
		if err != nil {
			log.WithError(err).Errorln("error with new auth config")
			return err
		}
		authConfig = ac.Configs["docker.io"]
	}

	err = client.PullImage(docker.PullImageOptions{Repository: repo, Tag: tag}, authConfig)
	if err != nil {
		log.WithError(err).Errorln("error pulling image")
		return err
	}
	return nil
}
