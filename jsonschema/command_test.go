package jsonschema_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/jsonschema"
)

const (
	name   = "test"
	schema = `{
  "$schema": "http://json-schema.org/draft/2020-12/schema"
}
`
)

func TestNew(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		cmd     *jsonschema.Command
		want    string
		wantErr bool
	}{
		{
			name: "the schema and nothing else",
			cmd:  &jsonschema.Command{Schema: []byte(schema), Name: name},
			want: schema,
		},
		{
			name: "a schema without a trailing newline gets one",
			cmd:  &jsonschema.Command{Schema: []byte(`{"$schema": "x"}`)},
			want: "{\"$schema\": \"x\"}\n",
		},
		{
			name:    "an empty schema is an error, not an empty file",
			cmd:     &jsonschema.Command{},
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
			buf := &bytes.Buffer{}
			cmd := jsonschema.New(d.cmd)
			cmd.SetOut(buf)
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs(nil)
			// cobra prints the usage of a failing command to the output, which would
			// otherwise be mixed into the schema the test reads.
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			err := cmd.ExecuteContext(t.Context())
			if d.wantErr {
				if err == nil {
					t.Fatal("an error must be returned")
				}
				if buf.Len() != 0 {
					t.Fatalf("nothing must be written: %s", buf.String())
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(d.want, buf.String()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

// TestNew_help checks that the help names the program, because the example in it is
// the only place that says how the command is run.
func TestNew_help(t *testing.T) {
	t.Parallel()
	cmd := jsonschema.New(&jsonschema.Command{Schema: []byte(schema), Name: name})
	if !strings.Contains(cmd.Long, "$ test json-schema > test.json") {
		t.Fatalf("the help must show the example: %s", cmd.Long)
	}
	cmd = jsonschema.New(nil)
	if strings.Contains(cmd.Long, "$ ") {
		t.Fatalf("a command without a name must not show an example: %s", cmd.Long)
	}
}

func TestWith(t *testing.T) {
	t.Parallel()
	root := &cobra.Command{Use: name, Short: "a test command"}
	root.CompletionOptions.DisableDefaultCmd = true
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"json-schema"})
	if err := jsonschema.With(root, []byte(schema)).Execute(); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(schema, buf.String()); diff != "" {
		t.Fatal(diff)
	}
}
