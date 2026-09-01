// Package vcmd implements the version command of CLIs built with spf13/cobra.
//
// It is the cobra counterpart of suzuki-shunsuke/urfave-cli-v3-util's vcmd package.
package vcmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// Command describes the version the command reports.
type Command struct {
	// Name is the name of the program, which is printed by Print so that the line
	// stays readable when several versions are pasted together.
	Name string
	// Version is the version of the program. An empty version is reported as
	// "unknown", which is what a binary built without version information, such as
	// one from 'go install', has.
	Version string
	// SHA is the commit the binary was built from, reported by 'version --json'.
	SHA string
	// AfterVersion, when set, is called after the version has been printed, by both
	// the command and Print. It is the hook for output that should accompany the
	// version wherever it is asked for.
	//
	// urfave/cli offered the process global cli.VersionPrinter for this, which had to
	// be saved and restored around every test that touched it; a field here has no
	// such reach.
	AfterVersion func()
}

// New returns the 'version' command, which prints the version on its own line, or as
// JSON with --json, so that it can be read by a script as well as by a person.
//
// The command writes to the output of the cobra command it runs as, so it follows
// whatever SetOut the root was given rather than needing a writer of its own.
func New(cmd *Command) *cobra.Command {
	if cmd == nil {
		cmd = &Command{}
	}
	asJSON := false
	c := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if asJSON {
				if err := cmd.printJSON(c.OutOrStdout()); err != nil {
					return err
				}
			} else {
				fmt.Fprintln(c.OutOrStdout(), versionOrUnknown(cmd.Version))
			}
			cmd.afterVersion()
			return nil
		},
	}
	c.Flags().BoolVarP(&asJSON, "json", "j", false, "Output version in JSON format")
	return c
}

// With adds the version command to cmd, taking the name and the version from cmd
// itself, and returns cmd.
func With(cmd *cobra.Command, sha string) *cobra.Command {
	cmd.AddCommand(New(&Command{
		Name:    cmd.Name(),
		Version: cmd.Version,
		SHA:     sha,
	}))
	return cmd
}

// Print writes the line --version and -v print, which names the program as well as
// the version, unlike the version command, whose output is meant to be consumed.
func (cmd *Command) Print(w io.Writer) {
	fmt.Fprintf(w, "%s version %s\n", cmd.Name, versionOrUnknown(cmd.Version))
	cmd.afterVersion()
}

func (cmd *Command) printJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(map[string]string{
		"name":    cmd.Name,
		"version": cmd.Version,
		"sha":     cmd.SHA,
	}); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}
	return nil
}

func (cmd *Command) afterVersion() {
	if cmd.AfterVersion != nil {
		cmd.AfterVersion()
	}
}

// versionOrUnknown keeps the output a single, parsable word for a binary built
// without version information, such as one from 'go install'.
func versionOrUnknown(version string) string {
	if version == "" {
		return "unknown"
	}
	return version
}
