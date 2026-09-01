# cobra-util

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/cobra-util/main/LICENSE) [![Go Reference](https://pkg.go.dev/badge/github.com/suzuki-shunsuke/cobra-util.svg)](https://pkg.go.dev/github.com/suzuki-shunsuke/cobra-util)

cobra-util is a Go library for [spf13/cobra](https://pkg.go.dev/github.com/spf13/cobra).
It is the cobra counterpart of [urfave-cli-v3-util](https://github.com/suzuki-shunsuke/urfave-cli-v3-util).

- [cobrautil](cobrautil): The entry point of the process, and environment variable sources for flags
- [helpall](helpall): Output the help message of all commands
- [vcmd](vcmd): Version Command
- [keyring/ghtoken](keyring/ghtoken): Manage a GitHub Access token by keyring

## How To Use

[Example](cmd/hello/main.go)

`cobrautil.Main` runs the CLI, logs the error it returns, and exits with its exit code.
`cobrautil.Command` wires the root command to the process environment and adds `--version`, the `version` command, and the hidden `help-all` command.

```go
func main() {
	cobrautil.Main("hello", "1.0.0", run)
}

func run(ctx context.Context, logger *slogutil.Logger, env *cobrautil.Env) error {
	cmd := cobrautil.Command(env, &cobra.Command{
		Use:   "hello",
		Short: "A new cli application",
	}, &cobrautil.Options{
		SHA: "abc123",
	})
	cmd.AddCommand(/* ... */)
	return cmd.ExecuteContext(ctx)
}
```

### Environment variables

urfave/cli lets a flag fall back to environment variables with `cli.EnvVars`, which pflag has no equivalent of.
`cobrautil.Envs` binds environment variables to a flag, and `cobrautil.Command` applies them in the root command's `PersistentPreRunE`.

```go
cmd.Flags().String("token", "", "GitHub Access token")
cobrautil.Envs(cmd.Flags(), "token", "HELLO_GITHUB_TOKEN", "GITHUB_TOKEN")
```

The value comes from the command line first, then from the environment variables in the order they were bound, then from the flag's default.
The names are appended to the flag's usage as `[$HELLO_GITHUB_TOKEN] [$GITHUB_TOKEN]`, so the help says where else the value can come from.

Note that cobra runs only the closest `PersistentPreRunE` in the chain, so a command that sets one of its own must call `cobrautil.ApplyEnvs` itself.
