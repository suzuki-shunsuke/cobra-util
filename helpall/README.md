# helpall

`helpall` is a Go Package to show the help of all commands of CLIs built with [spf13/cobra](https://pkg.go.dev/github.com/spf13/cobra).
This is useful if you want to put the usage of CLI built with cobra into the document.

## How To Use

[Example](../cmd/hello/main.go)

Using this library, you can add a command `help-all` showing the help of all commands.
[cobrautil.Command](../cobrautil) adds it for you, so you only need this package directly if you don't use `cobrautil.Command`.

```go
import (
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/helpall"
)

// helpall.With appends the "help-all" command to the given command.
rootCmd := helpall.With(&cobra.Command{
	Use: "hello",
}, nil)
rootCmd.Execute()
```

`help-all` command outputs the help message.
You can put it into the document.

```sh
go run ./cmd/hello help-all > hello.md
```

Example: [hello.md](../hello.md)

### Customize the command

The function `helpall.New()` returns a `*cobra.Command`. You can customize the returned value.

e.g. Change the command name

```go
cmd := helpall.New(nil)
cmd.Use = "help-markdown" // Change the command name
rootCmd.AddCommand(cmd)
```

By default, the command is hidden, so it isn't shown in the help message.
You can show the command by changing the `Hidden` field.

```go
cmd := helpall.New(nil)
cmd.Hidden = false // Show the help of help-all
rootCmd.AddCommand(cmd)
```
