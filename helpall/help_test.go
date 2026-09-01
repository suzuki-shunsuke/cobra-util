package helpall_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/helpall"
)

// newRoot returns a tree holding one documented command with a subcommand, one
// hidden command, and the help-all command itself, which is what the skip rules
// have to deal with.
func newRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "test",
		Short: "a test command",
	}
	// cobra adds a completion command of its own, whose help is long and is not what
	// these tests are about.
	root.CompletionOptions.DisableDefaultCmd = true
	foo := &cobra.Command{
		Use:   "foo",
		Short: "foo command",
		Args:  cobra.NoArgs,
		RunE:  func(*cobra.Command, []string) error { return nil },
	}
	foo.AddCommand(&cobra.Command{
		Use:   "sub",
		Short: "sub command",
		Args:  cobra.NoArgs,
		RunE:  func(*cobra.Command, []string) error { return nil },
	})
	root.AddCommand(foo, &cobra.Command{
		Use:    "secret",
		Short:  "secret command",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE:   func(*cobra.Command, []string) error { return nil },
	})
	return helpall.With(root, nil)
}

func run(t *testing.T) string {
	t.Helper()
	buf := &bytes.Buffer{}
	root := newRoot()
	root.SetOut(buf)
	root.SetArgs([]string{"help-all"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestHeadings(t *testing.T) {
	t.Parallel()
	got := []string(nil)
	for line := range strings.SplitSeq(run(t), "\n") {
		if strings.HasPrefix(line, "#") {
			got = append(got, line)
		}
	}
	// The root's own help opens the document without a heading, the commands worth
	// documenting are nested under it, and help, help-all and the hidden command are
	// nowhere to be seen.
	want := []string{
		"## test foo",
		"### test foo sub",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestHelpFlagIsDocumented(t *testing.T) {
	t.Parallel()
	// cobra adds --help to the command it runs, which is help-all, so without
	// InitDefaultHelpFlag the documented help would be the one place it is missing.
	if diff := cmp.Diff(3, strings.Count(run(t), "-h, --help")); diff != "" {
		t.Fatal(diff)
	}
}

func TestConsoleBlocks(t *testing.T) {
	t.Parallel()
	got := run(t)
	for _, want := range []string{
		"```console\n$ test --help\n",
		"```console\n$ test foo --help\n",
		"```console\n$ test foo sub --help\n",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("the output must contain %q:\n%s", want, got)
		}
	}
	if diff := cmp.Diff(6, strings.Count(got, "```")); diff != "" {
		t.Fatal(diff)
	}
}

func TestNewIsHidden(t *testing.T) {
	t.Parallel()
	cmd := helpall.New(nil)
	if !cmd.Hidden {
		t.Fatal("the help-all command must be hidden by default")
	}
	if diff := cmp.Diff("help-all", cmd.Name()); diff != "" {
		t.Fatal(diff)
	}
}
