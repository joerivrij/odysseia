package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/cmd/command"
	"github.com/spf13/cobra"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "archimedes",
		Short: "Deploy everything related to odysseia",
		Long: `Create and script everything odysseia related.
Allows you to parse words from a txt file,
build all container images`,
	}
)

func main() {
	glg.Info("\n  ____  ____      __  __ __  ____  ___ ___    ___  ___      ___  _____\n /    ||    \\    /  ]|  |  ||    ||   |   |  /  _]|   \\    /  _]/ ___/\n|  o  ||  D  )  /  / |  |  | |  | | _   _ | /  [_ |    \\  /  [_(   \\_ \n|     ||    /  /  /  |  _  | |  | |  \\_/  ||    _]|  D  ||    _]\\__  |\n|  _  ||    \\ /   \\_ |  |  | |  | |   |   ||   [_ |     ||   [_ /  \\ |\n|  |  ||  .  \\\\     ||  |  | |  | |   |   ||     ||     ||     |\\    |\n|__|__||__|\\_| \\____||__|__||____||___|___||_____||_____||_____| \\___|\n                                                                      \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"εὕρηκα\"")
	glg.Info("\"I found it!\"")
	glg.Info(strings.Repeat("~", 37))

	rootCmd.AddCommand(
		command.ParseListToWords(),
		command.CreateImages(),
	)

	err := rootCmd.Execute()
	if err != nil {
		glg.Error(err)
	}
}
