// Package cobrautil provides the pieces every CLI of ours needs on top of
// spf13/cobra: the process entry point, the wiring of a root command to that
// entry point, and environment variable sources for flags.
//
// It is the cobra counterpart of suzuki-shunsuke/urfave-cli-v3-util, and the
// version and help-all commands it adds live in the sibling packages vcmd and
// helpall, which can also be used on their own.
package cobrautil

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

// Env is the process environment handed to Run, so that commands read it from here
// rather than reaching for the globals in os and can be driven by a test.
type Env struct {
	Program string
	Version string
	Stdin   *os.File
	Stdout  *os.File
	Stderr  *os.File
	Getenv  func(string) string
	Args    []string
}

// Run is the entry point of the CLI, called by Main with the process environment.
type Run func(ctx context.Context, logger *slogutil.Logger, env *Env) error

// ErrSilent is an error with an empty message, which Main logs nothing for. Wrap it
// with ecerror.Wrap to exit with a code without reporting a failure of the CLI's own,
// as when a command propagates the exit code of a command it ran.
var ErrSilent = errors.New("")

// Main runs the CLI and exits the process with the exit code of the error it returns.
// args are extra attributes for the logger.
func Main(name, version string, run Run, args ...any) {
	if code := core(name, version, run, args...); code != 0 {
		os.Exit(code)
	}
}

// core is Main without the os.Exit, so that a test can assert on the exit code.
func core(name, version string, run Run, args ...any) int {
	version = getVersion(version, debug.ReadBuildInfo)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	logger := slogutil.New(&slogutil.InputNew{
		Name:    name,
		Version: version,
		Out:     os.Stderr,
		Attrs:   args,
	})
	env := &Env{
		Program: name,
		Version: version,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Getenv:  os.Getenv,
		Args:    os.Args,
	}
	if err := run(ctx, logger, env); err != nil {
		if err.Error() != "" {
			slogerr.WithError(logger.Logger, err).Error(name + " failed")
		}
		return ecerror.GetExitCode(err)
	}
	return 0
}

// getVersion returns version if it isn't empty. Otherwise it falls back to the
// module version the Go toolchain embeds in the binary, so that a binary built
// without -ldflags, such as one from 'go install', still reports its version.
func getVersion(version string, readBuildInfo func() (*debug.BuildInfo, bool)) string {
	if version != "" {
		return version
	}
	if info, ok := readBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "unknown"
}

// RunEFunc is a command body that is given the logger built by Main, on top of what
// cobra passes to a cobra.Command's RunE.
type RunEFunc func(ctx context.Context, cmd *cobra.Command, args []string, logger *slogutil.Logger) error

// RunE adapts fn to a cobra.Command's RunE, binding the logger to it. It is how the
// logger reaches a command body without being passed through a global or stored on
// every command struct.
//
// The context is the one cobra runs the command with, which Execute derives from the
// context Main hands to Run, so a command body sees the signal handling set up there.
func RunE(fn RunEFunc, logger *slogutil.Logger) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return fn(cmd.Context(), cmd, args, logger)
	}
}
