package vcmd_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/vcmd"
)

const (
	name    = "test"
	version = "1.0.0"
	sha     = "abc123"
)

func TestNew(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name string
		cmd  *vcmd.Command
		args []string
		want string
	}{
		{
			name: "the version on its own line",
			cmd:  &vcmd.Command{Name: name, Version: version, SHA: sha},
			want: "1.0.0\n",
		},
		{
			name: "an empty version is unknown",
			cmd:  &vcmd.Command{Name: name},
			want: "unknown\n",
		},
		{
			name: "a nil command still runs",
			want: "unknown\n",
		},
		{
			name: "JSON",
			cmd:  &vcmd.Command{Name: name, Version: version, SHA: sha},
			args: []string{"--json"},
			want: "{\n  \"name\": \"test\",\n  \"sha\": \"abc123\",\n  \"version\": \"1.0.0\"\n}\n",
		},
		{
			name: "JSON keeps an empty version empty, unlike the plain output",
			cmd:  &vcmd.Command{Name: name},
			args: []string{"-j"},
			want: "{\n  \"name\": \"test\",\n  \"sha\": \"\",\n  \"version\": \"\"\n}\n",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			after := 0
			if d.cmd != nil {
				d.cmd.AfterVersion = func() { after++ }
			}
			buf := &bytes.Buffer{}
			cmd := vcmd.New(d.cmd)
			cmd.SetOut(buf)
			cmd.SetArgs(d.args)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(d.want, buf.String()); diff != "" {
				t.Fatal(diff)
			}
			if d.cmd != nil {
				if diff := cmp.Diff(1, after); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}

func TestWith(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	root := &cobra.Command{Use: name, Version: version}
	root.SetOut(buf)
	root.SetArgs([]string{"version", "--json"})
	if err := vcmd.With(root, sha).Execute(); err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"name\": \"test\",\n  \"sha\": \"abc123\",\n  \"version\": \"1.0.0\"\n}\n"
	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Fatal(diff)
	}
}

func TestCommandPrint(t *testing.T) {
	t.Parallel()
	data := []struct {
		name string
		cmd  *vcmd.Command
		want string
	}{
		{
			name: "the program name and the version",
			cmd:  &vcmd.Command{Name: name, Version: version},
			want: "test version 1.0.0\n",
		},
		{
			name: "an empty version is unknown",
			cmd:  &vcmd.Command{Name: name},
			want: "test version unknown\n",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			after := 0
			d.cmd.AfterVersion = func() { after++ }
			buf := &bytes.Buffer{}
			d.cmd.Print(buf)
			if diff := cmp.Diff(d.want, buf.String()); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(1, after); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
