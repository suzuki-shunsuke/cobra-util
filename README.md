# cobra-util

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/cobra-util/main/LICENSE) [![Go Reference](https://pkg.go.dev/badge/github.com/suzuki-shunsuke/cobra-util.svg)](https://pkg.go.dev/github.com/suzuki-shunsuke/cobra-util)

cobra-util is a Go library for [spf13/cobra](https://pkg.go.dev/github.com/spf13/cobra).
It is the cobra counterpart of [urfave-cli-v3-util](https://github.com/suzuki-shunsuke/urfave-cli-v3-util).

- [cobrautil](cobrautil): The entry point of the process, and environment variable sources for flags
- [helpall](helpall): Output the help message of all commands
- [vcmd](vcmd): Version Command
- [jsonschema](jsonschema): Output the JSON Schema of the configuration file
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

### JSON Schema

`jsonschema.With` adds the `json-schema` command, which outputs the JSON Schema of the configuration file so that editors such as VSCode can complete the configuration file and warn about invalid settings.
The schema is passed as bytes, so embed the schema the CLI generates with `go:embed` and it always describes the configuration the running version accepts.

`With` adds the command to the command it is given and returns it, so it can wrap `cobrautil.Command`.

```go
//go:embed json-schema/hello.json
var schema []byte

cmd := jsonschema.With(cobrautil.Command(env, rootCmd, opts), schema)
```

It takes the name of the program from the command it adds to, which is what the example in the help says.
Use `jsonschema.New` instead to name the program yourself, or to add the command somewhere other than the root.

```go
cmd.AddCommand(jsonschema.New(&jsonschema.Command{
	Schema: schema,
	Name:   "hello",
}))
```

```console
$ hello json-schema > hello.json
```
