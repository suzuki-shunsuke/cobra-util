package cobrautil_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/cobrautil"
)

const (
	name       = "test"
	version    = "1.0.0"
	fromEnv    = "from-env"
	envToken   = "TOKEN"
	envGHToken = "GITHUB_TOKEN" //nolint:gosec // This is the name of an environment variable, not a credential.
)

// newEnv returns an Env with no files, so that a test can attach a buffer with
// SetOut: Command only calls SetOut when Env.Stdout is set.
func newEnv(args ...string) *cobrautil.Env {
	return &cobrautil.Env{
		Program: name,
		Version: version,
		Args:    append([]string{name}, args...),
		Getenv:  func(string) string { return "" },
	}
}

func TestCommandVersionFlag(t *testing.T) {
	t.Parallel()
	for _, flag := range []string{"--version", "-v"} {
		t.Run(flag, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			after := 0
			cmd := cobrautil.Command(newEnv(flag), &cobra.Command{Use: name}, &cobrautil.Options{
				AfterVersion: func() { after++ },
			})
			cmd.SetOut(buf)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(name+" version "+version+"\n", buf.String()); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(1, after); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestCommandVersionSubCommand(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	cmd := cobrautil.Command(newEnv("version", "--json"), &cobra.Command{Use: name}, &cobrautil.Options{
		SHA: "abc123",
	})
	cmd.SetOut(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	want := `{
  "name": "test",
  "sha": "abc123",
  "version": "1.0.0"
}
`
	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Fatal(diff)
	}
}

func TestCommandNoAction(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	cmd := cobrautil.Command(newEnv(), &cobra.Command{
		Use:   name,
		Short: "a test command",
	}, nil)
	cmd.SetOut(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Usage:") {
		t.Fatalf("a root command with no action must show its help: %s", buf.String())
	}
}

func TestCommandKeepsRunE(t *testing.T) {
	t.Parallel()
	called := false
	cmd := cobrautil.Command(newEnv(), &cobra.Command{
		Use: name,
		RunE: func(*cobra.Command, []string) error {
			called = true
			return nil
		},
	}, nil)
	cmd.SetOut(&bytes.Buffer{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("the RunE of the given command must be kept")
	}
}

func TestCommandAppliesEnvs(t *testing.T) {
	t.Parallel()
	token := ""
	root := &cobra.Command{
		Use: name,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			token, err = cmd.Flags().GetString("token")
			return err //nolint:wrapcheck
		},
	}
	root.Flags().String("token", "", "the token")
	cobrautil.Envs(root.Flags(), "token", "TEST_TOKEN")
	env := newEnv()
	env.Getenv = func(k string) string {
		if k == "TEST_TOKEN" {
			return fromEnv
		}
		return ""
	}
	cmd := cobrautil.Command(env, root, nil)
	cmd.SetOut(&bytes.Buffer{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(fromEnv, token); diff != "" {
		t.Fatal(diff)
	}
}
