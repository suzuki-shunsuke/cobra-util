# vcmd

vcmd is a Go Package to implement `version` command of CLIs built with [spf13/cobra](https://pkg.go.dev/github.com/spf13/cobra).

## How To Use

[Example](../cmd/hello/main.go)

[cobrautil.Command](../cobrautil) adds the command for you, so you only need this package directly if you don't use `cobrautil.Command`.

```go
import (
	"github.com/spf13/cobra"
	"github.com/suzuki-shunsuke/cobra-util/vcmd"
)

// vcmd.With appends the "version" command to the given command,
// taking the name and the version from it.
rootCmd := vcmd.With(&cobra.Command{
	Use:     "hello",
	Version: "1.0.0",
}, "abc123")
```

The command prints the version on its own line, or as JSON with `--json`.

```console
$ hello version
1.0.0
$ hello version --json
{
  "name": "hello",
  "sha": "abc123",
  "version": "1.0.0"
}
```
