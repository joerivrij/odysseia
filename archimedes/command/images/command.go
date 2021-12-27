package images

import (
	"github.com/odysseia/archimedes/command/images/command"
	"github.com/spf13/cobra"
)

func Manager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "work with images",
		Long:  `Allows you to create and parse images`,
	}

	cmd.AddCommand(
		command.CreateImages(),
	)

	return cmd
}
