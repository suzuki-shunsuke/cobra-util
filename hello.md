```console
$ hello --help
A new cli application

Usage:
  hello [flags]
  hello [command]

Available Commands:
  bar         bar command
  completion  Generate the autocompletion script for the specified shell
  foo         foo command
  help        Help about any command
  version     Show version

Flags:
  -h, --help      help for hello
  -v, --version   print the version

Use "hello [command] --help" for more information about a command.
```

## hello bar

```console
$ hello bar --help
This is a bar command

Usage:
  hello bar [flags]

Flags:
  -h, --help   help for bar
```

## hello completion

```console
$ hello completion --help
Generate the autocompletion script for hello for the specified shell.
See each sub-command's help for details on how to use the generated script.

Usage:
  hello completion [command]

Available Commands:
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh

Flags:
  -h, --help   help for completion

Use "hello completion [command] --help" for more information about a command.
```

### hello completion bash

```console
$ hello completion bash --help
Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(hello completion bash)

To load completions for every new session, execute once:

#### Linux:

	hello completion bash > /etc/bash_completion.d/hello

#### macOS:

	hello completion bash > $(brew --prefix)/etc/bash_completion.d/hello

You will need to start a new shell for this setup to take effect.

Usage:
  hello completion bash

Flags:
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### hello completion fish

```console
$ hello completion fish --help
Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	hello completion fish | source

To load completions for every new session, execute once:

	hello completion fish > ~/.config/fish/completions/hello.fish

You will need to start a new shell for this setup to take effect.

Usage:
  hello completion fish [flags]

Flags:
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### hello completion powershell

```console
$ hello completion powershell --help
Generate the autocompletion script for powershell.

To load completions in your current shell session:

	hello completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.

Usage:
  hello completion powershell [flags]

Flags:
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### hello completion zsh

```console
$ hello completion zsh --help
Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(hello completion zsh)

To load completions for every new session, execute once:

#### Linux:

	hello completion zsh > "${fpath[1]}/_hello"

#### macOS:

	hello completion zsh > $(brew --prefix)/share/zsh/site-functions/_hello

You will need to start a new shell for this setup to take effect.

Usage:
  hello completion zsh [flags]

Flags:
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

## hello foo

```console
$ hello foo --help
This is a foo command

Usage:
  hello foo [flags]

Flags:
  -h, --help   help for foo
```

## hello version

```console
$ hello version --help
Show version

Usage:
  hello version [flags]

Flags:
  -h, --help   help for version
  -j, --json   Output version in JSON format
```
