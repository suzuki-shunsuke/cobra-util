// Package jsonschema implements the json-schema command of CLIs built with
// spf13/cobra, which outputs the JSON Schema of the configuration file so that
// editors such as VSCode can complete the configuration file and warn about invalid
// settings.
//
// The schema is passed as bytes rather than read at run time, so that the CLI can
// embed it with go:embed and always output the schema of the configuration the
// running version accepts.
package jsonschema

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// errNoSchema is returned rather than writing nothing, because an empty output would
// otherwise be truncated into the configuration file the user redirected it to.
var errNoSchema = errors.New("no JSON Schema is embedded in this binary")

// Command describes the schema the command outputs.
type Command struct {
	// Schema is the JSON Schema, which is usually embedded with go:embed. The command
	// fails if it is empty.
	Schema []byte
	// Name is the name of the program, which appears in the example in the help so
	// that it reads as a command the user can run. It may be empty, in which case the
	// help drops the example; With takes it from the root command.
	Name string
}

// New returns the 'json-schema' command, which writes the schema to the output of the
// cobra command it runs as, so it follows whatever SetOut the root was given.
//
// Unlike the other commands of a CLI, it needs neither a logger nor the global flags:
// it reads no configuration and logs nothing.
func New(cmd *Command) *cobra.Command {
	if cmd == nil {
		cmd = &Command{}
	}
	return &cobra.Command{
		Use:   "json-schema",
		Short: "Output JSON Schema for the configuration file",
		Long:  cmd.long(),
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return cmd.Write(c.OutOrStdout())
		},
	}
}

// With adds the json-schema command to rootCmd, taking the name from rootCmd itself,
// and returns rootCmd.
func With(rootCmd *cobra.Command, schema []byte) *cobra.Command {
	rootCmd.AddCommand(New(&Command{
		Schema: schema,
		Name:   rootCmd.Name(),
	}))
	return rootCmd
}

// Write writes the schema to w, terminated by a newline. It is exported so that the
// schema can also be written somewhere else, such as to the file a test compares with
// the committed schema.
func (cmd *Command) Write(w io.Writer) error {
	if len(cmd.Schema) == 0 {
		return errNoSchema
	}
	if _, err := w.Write(cmd.Schema); err != nil {
		return fmt.Errorf("write the JSON Schema: %w", err)
	}
	if cmd.Schema[len(cmd.Schema)-1] != '\n' {
		// A generated schema usually ends with a newline already; one that doesn't
		// would leave the shell prompt on the same line as the last brace.
		if _, err := fmt.Fprintln(w); err != nil {
			return fmt.Errorf("write a newline: %w", err)
		}
	}
	return nil
}

const (
	// longWithName explains the command and shows how it is run, which needs the name
	// of the program because the command is only ever run through it.
	longWithName = `Output JSON Schema for the configuration file.

The schema is embedded in the binary, so it always describes the configuration
this version of %[1]s accepts. Editors such as VSCode use it to complete the
configuration file and to warn about invalid settings.

$ %[1]s json-schema > %[1]s.json`
	// longWithoutName is the help of a command built without a name, whose example
	// would otherwise name a command that doesn't exist.
	longWithoutName = `Output JSON Schema for the configuration file.

The schema is embedded in the binary, so it always describes the configuration
this version accepts. Editors such as VSCode use it to complete the
configuration file and to warn about invalid settings.`
)

func (cmd *Command) long() string {
	if cmd.Name == "" {
		return longWithoutName
	}
	return fmt.Sprintf(longWithName, cmd.Name)
}
