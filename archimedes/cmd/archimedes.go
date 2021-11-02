package main

import (
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/command/images"
	"github.com/odysseia/archimedes/command/kubernetes"
	"github.com/odysseia/archimedes/command/parse"
	"github.com/odysseia/archimedes/command/vault"
	"github.com/spf13/cobra"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "archimedes",
		Short: "Deploy everything related to odysseia",
		Long: `Create and script everything odysseia related.
Allows you to parse words from a txt file,
build all container images
work with vault and much more is coming`,
	}
)

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ARCHIMEDES
	glg.Info("\n  ____  ____      __  __ __  ____  ___ ___    ___  ___      ___  _____\n /    ||    \\    /  ]|  |  ||    ||   |   |  /  _]|   \\    /  _]/ ___/\n|  o  ||  D  )  /  / |  |  | |  | | _   _ | /  [_ |    \\  /  [_(   \\_ \n|     ||    /  /  /  |  _  | |  | |  \\_/  ||    _]|  D  ||    _]\\__  |\n|  _  ||    \\ /   \\_ |  |  | |  | |   |   ||   [_ |     ||   [_ /  \\ |\n|  |  ||  .  \\\\     ||  |  | |  | |   |   ||     ||     ||     |\\    |\n|__|__||__|\\_| \\____||__|__||____||___|___||_____||_____||_____| \\___|\n                                                                      \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"εὕρηκα\"")
	glg.Info("\"I found it!\"")
	glg.Info(strings.Repeat("~", 37))

	rootCmd.AddCommand(
		images.Manager(),
		parse.Manager(),
		vault.Manager(),
		kubernetes.Manager(),
	)

	err := rootCmd.Execute()
	if err != nil {
		glg.Error(err)
	}
}
