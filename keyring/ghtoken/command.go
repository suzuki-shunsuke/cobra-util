// Package ghtoken implements the commands managing a GitHub Access token in the
// keyring of the operating system, and the oauth2 token source reading it back.
//
// It is the cobra counterpart of suzuki-shunsuke/urfave-cli-v3-util's
// keyring/ghtoken package.
package ghtoken

import (
	"github.com/spf13/cobra"
)

// InputSet holds the flags of the 'set' command.
type InputSet struct {
	Stdin bool
}

// Command returns the 'token' command and its subcommands, which manage the GitHub
// Access token in the keyring.
func Command(actor *Actor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage GitHub Access token",
		Long:  `Manage GitHub Access token by keyring.`,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(setCommand(actor), removeCommand(actor))
	return cmd
}

func setCommand(actor *Actor) *cobra.Command {
	input := &InputSet{}
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set GitHub Access token",
		Long:  `Set GitHub Access token to keyring.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return actor.Set(cmd.Context(), input)
		},
	}
	cmd.Flags().BoolVar(&input.Stdin, "stdin", false, "Read GitHub Access token from stdin")
	return cmd
}

func removeCommand(actor *Actor) *cobra.Command {
	return &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove GitHub Access token",
		Long:    `Remove GitHub Access token from keyring.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return actor.Remove()
		},
	}
}
