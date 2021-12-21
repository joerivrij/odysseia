package kubernetes

import (
	"github.com/odysseia/archimedes/command/kubernetes/command"
	"github.com/spf13/cobra"
)

func Manager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "work with your kubernetes cluster",
		Long:  `Allows you to operate the clustered kubernetes using a local client`,
	}

	cmd.AddCommand(
		command.CreateSecret(),
		command.CreateElasticCerts(),
	)

	return cmd
}
