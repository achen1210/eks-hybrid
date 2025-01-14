package config

import (
	"github.com/aws/eks-hybrid/internal/cli"
)

const configHelpText = `Examples:
  # Check configuration file
  nodeadm config check --config nodeConfig.yaml

Documentation:
  https://docs.aws.amazon.com/eks/latest/userguide/hybrid-nodes-prerequisites.html`

func NewConfigCommand() cli.Command {
	container := cli.NewCommandContainer("config", "Manage and validate hybrid node configuration.")
	container.AddCommand(NewCheckCommand())
    container.Flaggy().AdditionalHelpAppend = configHelpText
	return container.AsCommand()
}
