package vault

import (
	"github.com/odysseia/archimedes/command/vault/command"
	"github.com/spf13/cobra"
)

func Manager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "work with the clustered vault",
		Long:  `Allows you to operate the clustered vault using a local client`,
	}

	cmd.AddCommand(
		command.Auth(),
		command.Unseal(),
		command.Policy(),
		command.Init(),
		command.New(),
		command.TLS(),
	)

	return cmd
}
