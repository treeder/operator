/*
This is a special command, not intended for user's use, operator uses it to pull private images on remote machine.
*/

package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var (
	image    string
	username string
	password string
)

var pullCmd = &cobra.Command{
	Use:   "pull [host]",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Infoln("args", args)
		if len(args) == 0 {
			log.Fatalln("Must have image")
		}
		image = args[0]

		endpoint := "unix:///var/run/docker.sock"
		client, err := docker.NewClient(endpoint)
		if err != nil {
			log.WithError(err).Fatalln("Error making docker client!")
			return
		}
		log.Infoln("username:", username, "password:", password)

		repo, tag := docker.ParseRepositoryTag(image)

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
				return
			}
			authConfig = ac.Configs["docker.io"]
		}

		err = client.PullImage(docker.PullImageOptions{Repository: repo, Tag: tag}, authConfig)
		if err != nil {
			log.WithError(err).Errorln("error pulling image")
			return
		}

	},
}

func init() {
	RootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringVarP(&image, "image", "i", "", "Docker hub image")
	pullCmd.Flags().StringVarP(&username, "username", "u", "", "Docker hub username")
	pullCmd.Flags().StringVarP(&password, "password", "p", "", "Docker hub password")
}
