package secrets

import (
	"github.com/PrivixAI-labs/Privix-node/command/helper"
	"github.com/PrivixAI-labs/Privix-node/command/secrets/authtoken"
	"github.com/PrivixAI-labs/Privix-node/command/secrets/generate"
	initCmd "github.com/PrivixAI-labs/Privix-node/command/secrets/init"
	"github.com/PrivixAI-labs/Privix-node/command/secrets/output"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	secretsCmd := &cobra.Command{
		Use:   "secrets",
		Short: "Top level SecretsManager command for interacting with secrets functionality. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(secretsCmd)

	registerSubcommands(secretsCmd)

	return secretsCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// secrets init
		initCmd.GetCommand(),
		// secrets generate
		generate.GetCommand(),
		// secrets output public data
		output.GetCommand(),
		// secrets authtoken
		authtoken.GetCommand(),
	)
}
