package ghtoken

import (
	"context"
	"fmt"
	"io"
	"syscall"

	"github.com/charmbracelet/x/term"
)

type PasswordReader struct {
	stdout io.Writer
}

func NewPasswordReader(stdout io.Writer) *PasswordReader {
	return &PasswordReader{
		stdout: stdout,
	}
}

func (p *PasswordReader) ReadPassword(ctx context.Context) ([]byte, error) {
	fmt.Fprint(p.stdout, "Enter a GitHub access token: ")
	type result struct {
		data []byte
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		b, err := term.ReadPassword(uintptr(syscall.Stdin))
		ch <- result{b, err}
	}()
	select {
	case <-ctx.Done():
		fmt.Fprintln(p.stdout, "")
		return nil, ctx.Err()
	case r := <-ch:
		fmt.Fprintln(p.stdout, "")
		if r.err != nil {
			return nil, fmt.Errorf("read a GitHub access token from terminal: %w", r.err)
		}
		return r.data, nil
	}
}
