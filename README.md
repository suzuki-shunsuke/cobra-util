# cobra-util

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/cobra-util/main/LICENSE) [![Go Reference](https://pkg.go.dev/badge/github.com/suzuki-shunsuke/cobra-util.svg)](https://pkg.go.dev/github.com/suzuki-shunsuke/cobra-util)

cobra-util is a Go library for [spf13/cobra](https://pkg.go.dev/github.com/spf13/cobra).
It is the cobra counterpart of [urfave-cli-v3-util](https://github.com/suzuki-shunsuke/urfave-cli-v3-util).

- [cobrautil](cobrautil): The entry point of the process, and environment variable sources for flags
- [helpall](helpall): Output the help message of all commands
- [vcmd](vcmd): Version Command
- [jsonschema](jsonschema): Output the JSON Schema of the configuration file
- [docs](docs): Output the embedded documents for coding agents
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

### Documents for coding agents

`docs.With` adds the `docs` command, which lists and outputs the documents embedded in the binary.
Coding agents read them before answering questions about the CLI or troubleshooting its errors, so they always see the documents of the version that is running rather than whatever a web search turns up.

The documents are Markdown files with a YAML frontmatter holding their description, embedded with `go:embed`.

```markdown
---
description: How to install hello. Use to pick an installation method and to verify the binary.
---

# Install
```

```go
//go:embed docs/*.md
var docsFS embed.FS

sub, err := fs.Sub(docsFS, "docs")
if err != nil {
	return err
}
cmd := docs.With(cobrautil.Command(env, rootCmd, opts), sub)
```

`docs list` outputs the name and the description of every document as JSON, and `docs show <name>` outputs one of them.
A name that doesn't exist is reported with the names that do, so an agent that guessed wrong recovers without listing the documents again.
Documents in subdirectories are served too, named by their path, such as `docs show codes/001`, so documentation grouped in directories doesn't have to be flattened.

```console
$ hello docs list
{
  "results": [
    {
      "name": "install",
      "description": "How to install hello. Use to pick an installation method and to verify the binary."
    }
  ],
  "help": "Run `hello docs show <name>` to see the details of each document."
}

$ hello docs show install
```

Use `docs.New` instead of `docs.With` to name the program yourself, or to add the command somewhere other than the root.

```go
docsCmd := &docs.Command{
	FS:   sub,
	Name: "hello",
}
cmd.AddCommand(docs.New(docsCmd))
```

An agent only learns that these documents exist if a command it runs says so, so log `Command.Hint` from the commands it is likely to run.
Log it to stderr rather than stdout, so that it doesn't break scripts that parse the output, and at the info level so that a lower log level silences it.

```go
logger.Info(docsCmd.Hint())
```
