package docs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/spf13/cobra"
)

// Command describes the documents the command lists and outputs.
type Command struct {
	// FS holds the documents as Markdown files with a YAML frontmatter, and is
	// usually embedded with go:embed. The command fails if it is nil.
	//
	// Names are paths relative to the root of FS, so an embed.FS holding the
	// documents in a subdirectory has to be narrowed with fs.Sub first, and a
	// document in a subdirectory of the documentation is named by its path, such as
	// "codes/001".
	FS fs.FS
	// Name is the name of the program, which appears in the help, in the output of
	// 'docs list', and in Hint, so that they read as commands that can be run. It may
	// be empty, in which case they name no program; With takes it from the root
	// command.
	Name string
}

// New returns the 'docs' command and its 'list' and 'show' subcommands, which write
// to the output of the cobra command they run as, so they follow whatever SetOut the
// root was given.
//
// Unlike the other commands of a CLI, they need neither a logger nor the global
// flags: they read no configuration and log nothing.
func New(cmd *Command) *cobra.Command {
	if cmd == nil {
		cmd = &Command{}
	}
	c := &cobra.Command{
		Use:   "docs",
		Short: "Output documents for coding agents",
		Long: `List and output the documents embedded in the binary.
This command is for coding agents to read the documents and solve problems.`,
		Args: cobra.NoArgs,
	}
	c.AddCommand(cmd.listCommand(), cmd.showCommand())
	return c
}

// With adds the docs command to rootCmd, taking the name from rootCmd itself, and
// returns rootCmd.
func With(rootCmd *cobra.Command, fsys fs.FS) *cobra.Command {
	rootCmd.AddCommand(New(&Command{
		FS:   fsys,
		Name: rootCmd.Name(),
	}))
	return rootCmd
}

// Hint points coding agents at the docs command. It is here, next to the command it
// advertises, so that several commands can emit it: an agent that never makes the
// CLI fail only learns that these documents exist if a command it does run says so.
//
// Log it to stderr rather than writing it to stdout, so that it doesn't break
// scripts that parse the output, and at the info level so that '--log-level warn'
// silences it.
func (cmd *Command) Hint() string {
	name := cmd.Name
	subject := "this CLI"
	if name != "" {
		subject = name
		name += " "
	}
	return fmt.Sprintf("If you are a coding agent, run `%[1]sdocs list` to list the documentation and `%[1]sdocs show <name>` to read it before answering questions about %[2]s or troubleshooting its errors.", name, subject)
}

// Names lists the names of the documents, including the ones in subdirectories,
// which are the names 'docs show' takes.
func (cmd *Command) Names() ([]string, error) {
	if cmd.FS == nil {
		return nil, errNoDocs
	}
	return Names(cmd.FS)
}

// results is the output of 'docs list'. It is JSON rather than a table because it is
// read by a coding agent, and it carries the help so that the agent learns how to
// read a document from the list itself.
type results struct {
	Results []*Result `json:"results"`
	Help    string    `json:"help"`
}

// List writes the name and the description of every document to w as JSON.
func (cmd *Command) List(w io.Writer) error {
	names, err := cmd.Names()
	if err != nil {
		return fmt.Errorf("list documents: %w", err)
	}
	rs := make([]*Result, 0, len(names))
	for _, name := range names {
		b, err := fs.ReadFile(cmd.FS, name+Ext)
		if err != nil {
			return fmt.Errorf("read a document file: %w", err)
		}
		result := &Result{
			Name: name,
		}
		if err := Parse(b, result); err != nil {
			return fmt.Errorf("parse the frontmatter of the document %s: %w", name, err)
		}
		rs = append(rs, result)
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	// The output is read rather than embedded in a web page, and escaping would turn
	// the placeholder in the help and any "<" in a description into \u003c.
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(&results{
		Results: rs,
		Help:    fmt.Sprintf("Run `%sdocs show <name>` to see the details of each document.", cmd.prefix()),
	}); err != nil {
		return fmt.Errorf("encode documents as JSON: %w", err)
	}
	return nil
}

// Show writes the content of the document named name to w. The name is the path of
// the document in FS without the extension, which is what Names lists.
func (cmd *Command) Show(w io.Writer, name string) error {
	if cmd.FS == nil {
		return errNoDocs
	}
	if name == "" {
		return cmd.withAvailable(errors.New("a document name is required"))
	}
	b, err := fs.ReadFile(cmd.FS, name+Ext)
	if err != nil {
		return cmd.withAvailable(fmt.Errorf("the document %s isn't found", name))
	}
	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write the document: %w", err)
	}
	if len(b) > 0 && b[len(b)-1] != '\n' {
		// A Markdown file usually ends with a newline already; one that doesn't would
		// leave the shell prompt on the same line as the last line of the document.
		if _, err := fmt.Fprintln(w); err != nil {
			return fmt.Errorf("write a newline: %w", err)
		}
	}
	return nil
}

// withAvailable appends the available document names to err, so that a coding agent
// that guessed a wrong name can recover without running 'docs list' again.
func (cmd *Command) withAvailable(err error) error {
	names, nameErr := cmd.Names()
	if nameErr != nil {
		return fmt.Errorf("%w: %w", err, nameErr)
	}
	return fmt.Errorf("%w. Available documents: %s", err, strings.Join(names, ", "))
}

// prefix is the program name as it is written before a subcommand, which is nothing
// at all for a command built without a name.
func (cmd *Command) prefix() string {
	if cmd.Name == "" {
		return ""
	}
	return cmd.Name + " "
}

func (cmd *Command) listCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the name and the description of every document",
		Long: fmt.Sprintf(`List the name and the description of every document as JSON.
The name is what "%sdocs show" takes, and the description says what the document
covers, so that only the documents worth reading are read.`, cmd.prefix()),
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return cmd.List(c.OutOrStdout())
		},
	}
}

func (cmd *Command) showCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [<doc>]",
		Short: "Output the content of a given document",
		Long: fmt.Sprintf(`Output the content of a given document.
This command needs a document name.
To see the names, list the documents with "%sdocs list".`, cmd.prefix()),
		// The name is optional so that running the command without one reports the
		// available documents instead of cobra's message about the argument count.
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmd.completeNames,
		RunE: func(c *cobra.Command, positional []string) error {
			name := ""
			if len(positional) > 0 {
				name = positional[0]
			}
			return cmd.Show(c.OutOrStdout(), name)
		},
	}
}

// completeNames completes the document name, which is otherwise only known by
// running 'docs list'.
func (cmd *Command) completeNames(_ *cobra.Command, positional []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(positional) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names, err := cmd.Names()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
