package settoken

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func (c *Controller) Set(ctx context.Context, logger *slog.Logger) error {
	text, err := c.get(ctx, logger)
	if err != nil {
		return fmt.Errorf("get a GitHub access Token: %w", err)
	}
	if err := c.tokenManager.Set(strings.TrimSpace(string(text))); err != nil {
		return fmt.Errorf("set a GitHub access Token to the secret store: %w", err)
	}
	return nil
}

func (c *Controller) get(ctx context.Context, logger *slog.Logger) ([]byte, error) {
	if c.param.IsStdin {
		type result struct {
			data []byte
			err  error
		}
		ch := make(chan result, 1)
		go func() {
			s, err := io.ReadAll(c.param.Stdin)
			ch <- result{s, err}
		}()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case r := <-ch:
			if r.err != nil {
				return nil, fmt.Errorf("read a GitHub access token from stdin: %w", r.err)
			}
			logger.Debug("read a GitHub access token from stdin")
			return r.data, nil
		}
	}
	text, err := c.term.ReadPassword(ctx)
	if err != nil {
		return nil, fmt.Errorf("read a GitHub access Token from terminal: %w", err)
	}
	return text, nil
}
