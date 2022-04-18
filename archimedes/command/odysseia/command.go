package odysseia

import (
	"github.com/odysseia/archimedes/command/odysseia/command"
	"github.com/spf13/cobra"
)

func Manager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "odysseia",
		Short: "meta odysseia commands",
		Long:  `Install odysseia from scratch, currently only tested and working on docker desktop (macos)`,
	}

	cmd.AddCommand(
		command.Install(),
		command.Setup(),
	)

	return cmd
}
