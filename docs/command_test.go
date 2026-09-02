package docs_test

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/docs"
)

const (
	name    = "test"
	install = "install"
	show    = "show"
)

// testFS is a directory of documents as a CLI embeds it: Markdown files with a
// frontmatter, next to a file that isn't a document at all.
func testFS() fs.FS {
	return fstest.MapFS{
		"install.md": &fstest.MapFile{Data: []byte("---\ndescription: How to install.\n---\n\n# Install\n")},
		"usage.md":   &fstest.MapFile{Data: []byte("---\ndescription: How to use.\n---\n\n# Usage\n")},
		"logo.png":   &fstest.MapFile{Data: []byte("not a document")},
		// A document in a subdirectory is named by its path, so that a CLI whose
		// documentation is grouped in directories serves all of it.
		"codes/001.md": &fstest.MapFile{Data: []byte("---\ndescription: What error 001 means.\n---\n\n# 001\n")},
	}
}

func run(t *testing.T, cmd *docs.Command, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	c := docs.New(cmd)
	c.SetOut(buf)
	c.SetErr(&bytes.Buffer{})
	c.SetArgs(args)
	// cobra prints the usage of a failing command to the output, which would
	// otherwise be mixed into the document the test reads.
	c.SilenceUsage = true
	c.SilenceErrors = true
	err := c.ExecuteContext(t.Context())
	return buf.String(), err
}

func TestNew_list(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		cmd     *docs.Command
		want    string
		wantErr bool
	}{
		{
			name: "the documents and how to read them",
			cmd:  &docs.Command{FS: testFS(), Name: name},
			want: `{
  "results": [
    {
      "name": "codes/001",
      "description": "What error 001 means."
    },
    {
      "name": "install",
      "description": "How to install."
    },
    {
      "name": "usage",
      "description": "How to use."
    }
  ],
  "help": "Run ` + "`test docs show <name>`" + ` to see the details of each document."
}
`,
		},
		{
			name:    "a document without a frontmatter is a bug of the CLI",
			cmd:     &docs.Command{FS: fstest.MapFS{"broken.md": &fstest.MapFile{Data: []byte("# Broken\n")}}},
			wantErr: true,
		},
		{
			name:    "no documents at all is an error, not an empty list",
			cmd:     &docs.Command{},
			wantErr: true,
		},
		{
			name:    "a nil command still runs",
			wantErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			out, err := run(t, d.cmd, "list")
			if d.wantErr {
				if err == nil {
					t.Fatal("an error must be returned")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(d.want, out); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestNew_show(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		cmd     *docs.Command
		args    []string
		want    string
		wantErr string
	}{
		{
			name: "the document and nothing else",
			cmd:  &docs.Command{FS: testFS(), Name: name},
			args: []string{show, install},
			want: "---\ndescription: How to install.\n---\n\n# Install\n",
		},
		{
			name: "a document in a subdirectory is named by its path",
			cmd:  &docs.Command{FS: testFS(), Name: name},
			args: []string{show, "codes/001"},
			want: "---\ndescription: What error 001 means.\n---\n\n# 001\n",
		},
		{
			name:    "a name that escapes the documents is not found",
			cmd:     &docs.Command{FS: testFS(), Name: name},
			args:    []string{show, "../secret"},
			wantErr: "the document ../secret isn't found",
		},
		{
			name: "a document without a trailing newline gets one",
			cmd:  &docs.Command{FS: fstest.MapFS{"a.md": &fstest.MapFile{Data: []byte("# A")}}},
			args: []string{show, "a"},
			want: "# A\n",
		},
		{
			name:    "a wrong name reports the available documents",
			cmd:     &docs.Command{FS: testFS(), Name: name},
			args:    []string{show, "instal"},
			wantErr: "Available documents: codes/001, install, usage",
		},
		{
			name:    "no name reports the available documents too",
			cmd:     &docs.Command{FS: testFS(), Name: name},
			args:    []string{show},
			wantErr: "Available documents: codes/001, install, usage",
		},
		{
			name:    "no documents at all",
			cmd:     &docs.Command{},
			args:    []string{show, install},
			wantErr: "no documents are embedded in this binary",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			out, err := run(t, d.cmd, d.args...)
			if d.wantErr != "" {
				if err == nil {
					t.Fatal("an error must be returned")
				}
				if !strings.Contains(err.Error(), d.wantErr) {
					t.Fatalf("the error must contain %q: %v", d.wantErr, err)
				}
				if out != "" {
					t.Fatalf("nothing must be written: %s", out)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(d.want, out); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

// TestCommand_Hint checks that the hint names the program, because it is the only
// place that says how the docs command is run to the agent that reads it.
func TestCommand_Hint(t *testing.T) {
	t.Parallel()
	cmd := &docs.Command{FS: testFS(), Name: name}
	if !strings.Contains(cmd.Hint(), "`test docs list`") {
		t.Fatalf("the hint must show the command: %s", cmd.Hint())
	}
	cmd = &docs.Command{FS: testFS()}
	if !strings.Contains(cmd.Hint(), "`docs list`") || strings.Contains(cmd.Hint(), "  ") {
		t.Fatalf("a command without a name must not name a program: %s", cmd.Hint())
	}
}

func TestCommand_Names(t *testing.T) {
	t.Parallel()
	cmd := &docs.Command{FS: testFS()}
	names, err := cmd.Names()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{"codes/001", "install", "usage"}, names); diff != "" {
		t.Fatal(diff)
	}
	if _, err := (&docs.Command{}).Names(); err == nil {
		t.Fatal("an error must be returned when no documents are embedded")
	}
}

func TestWith(t *testing.T) {
	t.Parallel()
	root := &cobra.Command{Use: name, Short: "a test command"}
	root.CompletionOptions.DisableDefaultCmd = true
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"docs", show, "usage"})
	if err := docs.With(root, testFS()).Execute(); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("---\ndescription: How to use.\n---\n\n# Usage\n", buf.String()); diff != "" {
		t.Fatal(diff)
	}
}
