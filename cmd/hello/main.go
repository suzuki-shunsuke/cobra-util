// Package main is an example CLI built with cobra-util. Its help is what hello.md
// holds, which is regenerated with 'go run ./cmd/hello help-all > hello.md'.
package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/cobrautil"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

const (
	name    = "hello"
	version = "1.0.0"
	sha     = "abc123"
)

func main() {
	cobrautil.Main(name, version, run)
}

func run(ctx context.Context, _ *slogutil.Logger, env *cobrautil.Env) error {
	cmd := cobrautil.Command(env, &cobra.Command{
		Use:   name,
		Short: "A new cli application",
	}, &cobrautil.Options{
		SHA: sha,
	})
	cmd.AddCommand(
		&cobra.Command{
			Use:   "foo",
			Short: "foo command",
			Long:  "This is a foo command",
			Args:  cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				return nil
			},
		},
		&cobra.Command{
			Use:   "bar",
			Short: "bar command",
			Long:  "This is a bar command",
			Args:  cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				return nil
			},
		},
	)
	return cmd.ExecuteContext(ctx) //nolint:wrapcheck
}
