package config

import (
	"github.com/aws/eks-hybrid/internal/cli"
)

func NewConfigCommand() cli.Command {
	container := cli.NewCommandContainer("config", "Manage configuration")
	container.AddCommand(NewCheckCommand())
	container.Flaggy().AdditionalHelpAppend = cli.DocsLink
	return container.AsCommand()
}
