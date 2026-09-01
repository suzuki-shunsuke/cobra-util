package cobrautil_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/suzuki-shunsuke/cobra-util/cobrautil"
)

// fallback is the value of the second environment variable bound to a flag, which
// is applied only when the first one has none.
const fallback = "second"

func TestEnvs(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name  string
		envs  []string
		usage string
	}{
		{
			name:  "no environment variable leaves the usage alone",
			usage: "the token",
		},
		{
			name:  "one environment variable",
			envs:  []string{envToken},
			usage: "the token [$TOKEN]",
		},
		{
			name:  "several environment variables are appended in order",
			envs:  []string{envToken, envGHToken},
			usage: "the token [$TOKEN] [$GITHUB_TOKEN]",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			fs.String("token", "", "the token")
			cobrautil.Envs(fs, "token", d.envs...)
			if diff := cmp.Diff(d.usage, fs.Lookup("token").Usage); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestEnvsUnknownFlag(t *testing.T) {
	t.Parallel()
	defer func() {
		if recover() == nil {
			t.Fatal("Envs must panic for a flag that is not registered")
		}
	}()
	cobrautil.Envs(pflag.NewFlagSet("test", pflag.ContinueOnError), "token", "TOKEN")
}

func TestApplyEnvs(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name   string
		args   []string
		env    map[string]string
		getenv bool
		token  string
		isErr  bool
	}{
		{
			name:   "the default is kept when nothing sets the flag",
			getenv: true,
			token:  "default",
		},
		{
			name:   "the environment variable sets the flag",
			env:    map[string]string{envToken: fromEnv},
			getenv: true,
			token:  fromEnv,
		},
		{
			name:   "the command line beats the environment variable",
			args:   []string{"--token", "from-args"},
			env:    map[string]string{envToken: fromEnv},
			getenv: true,
			token:  "from-args",
		},
		{
			name:   "the first environment variable with a value wins",
			env:    map[string]string{envToken: "first", envGHToken: fallback},
			getenv: true,
			token:  "first",
		},
		{
			name:   "an empty value is treated as unset",
			env:    map[string]string{envToken: "", envGHToken: fallback},
			getenv: true,
			token:  fallback,
		},
		{
			name:   "a nil getenv applies nothing",
			env:    map[string]string{envToken: fromEnv},
			getenv: false,
			token:  "default",
		},
		{
			name:   "a value the flag can't parse is an error",
			env:    map[string]string{"VERBOSE": "yes"},
			getenv: true,
			isErr:  true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			cmd := &cobra.Command{Use: name}
			cmd.Flags().String("token", "default", "the token")
			cmd.Flags().Bool("verbose", false, "verbose")
			cobrautil.Envs(cmd.Flags(), "token", envToken, envGHToken)
			cobrautil.Envs(cmd.Flags(), "verbose", "VERBOSE")
			if err := cmd.Flags().Parse(d.args); err != nil {
				t.Fatal(err)
			}
			var getenv func(string) string
			if d.getenv {
				getenv = func(k string) string { return d.env[k] }
			}
			err := cobrautil.ApplyEnvs(cmd, getenv)
			if d.isErr {
				if err == nil {
					t.Fatal("error must be returned")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(d.token, token); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
