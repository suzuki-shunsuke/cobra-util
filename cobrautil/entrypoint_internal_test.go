package cobrautil

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

const (
	name    = "test"
	version = "1.0.0"
)

func TestCore(t *testing.T) {
	t.Parallel()
	data := []struct {
		name string
		err  error
		code int
	}{
		{
			name: "success",
			code: 0,
		},
		{
			name: "a plain error exits with 1",
			err:  errors.New("failed"),
			code: 1,
		},
		{
			name: "the exit code of the error is used",
			err:  ecerror.Wrap(errors.New("failed"), 3),
			code: 3,
		},
		{
			name: "a silent error is logged as nothing but still exits",
			err:  ecerror.Wrap(ErrSilent, 2),
			code: 2,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			called := false
			run := func(_ context.Context, logger *slogutil.Logger, env *Env) error {
				called = true
				if logger == nil {
					t.Fatal("a logger must be passed to Run")
				}
				if diff := cmp.Diff(name, env.Program); diff != "" {
					t.Fatal(diff)
				}
				if diff := cmp.Diff(version, env.Version); diff != "" {
					t.Fatal(diff)
				}
				if env.Getenv == nil {
					t.Fatal("Getenv must be set")
				}
				return d.err
			}
			if diff := cmp.Diff(d.code, core(name, version, run)); diff != "" {
				t.Fatal(diff)
			}
			if !called {
				t.Fatal("Run must be called")
			}
		})
	}
}

func TestRunE(t *testing.T) {
	t.Parallel()
	logger := slogutil.New(&slogutil.InputNew{Name: name})
	gotArgs := []string(nil)
	gotLogger := (*slogutil.Logger)(nil)
	cmd := &cobra.Command{
		Use:  name,
		Args: cobra.ArbitraryArgs,
		RunE: RunE(func(ctx context.Context, _ *cobra.Command, args []string, l *slogutil.Logger) error {
			if ctx == nil {
				t.Error("a context must be passed")
			}
			gotArgs = args
			gotLogger = l
			return nil
		}, logger),
	}
	cmd.SetArgs([]string{"foo", "bar"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{"foo", "bar"}, gotArgs); diff != "" {
		t.Fatal(diff)
	}
	if gotLogger != logger {
		t.Fatal("the logger must be passed through")
	}
}
