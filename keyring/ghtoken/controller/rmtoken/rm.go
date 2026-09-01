package rmtoken

import (
	"fmt"
	"log/slog"
)

func (c *Controller) Remove(logger *slog.Logger) error {
	if err := c.tokenManager.Remove(logger); err != nil {
		return fmt.Errorf("remove a GitHub access Token from the secret store: %w", err)
	}
	return nil
}
