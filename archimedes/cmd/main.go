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
		Long: `Deploy everything related to odysseia.
Allows you to parse words from a txt file`,
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
	)

	err := rootCmd.Execute()
	if err != nil {
		glg.Error(err)
	}
}
