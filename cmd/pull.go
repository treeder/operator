/*
This is a special command, not intended for user's use, operator uses it to pull private images on remote machine.
*/

package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/treeder/operator/commands"
	"golang.org/x/net/context"
)

var (
	username string
	password string
)

var pullCmd = &cobra.Command{
	Use:   "pull [host]",
	Short: "Pulls an image using docker api, NOT LIKE THE REST",
	Long: `This one is different than the rest, it is intended for use on a remote server to pull
	a private image using credentials. It does not work like the rest where it will apply to all servers in a set.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		log.Infoln("args", args)
		if len(args) == 0 {
			log.Fatalln("Must have image")
		}
		image := args[0]

		commands.Pull(ctx, image, username, password)

	},
}

func init() {
	RootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringVarP(&username, "username", "u", "", "Docker hub username")
	pullCmd.Flags().StringVarP(&password, "password", "p", "", "Docker hub password")
}
