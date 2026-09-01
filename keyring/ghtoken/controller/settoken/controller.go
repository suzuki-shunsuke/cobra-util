package settoken

import (
	"context"
	"io"
)

type Controller struct {
	param        *Param
	term         Terminal
	tokenManager TokenManager
}

func New(param *Param, term Terminal, tokenManager TokenManager) *Controller {
	return &Controller{
		param:        param,
		term:         term,
		tokenManager: tokenManager,
	}
}

type Param struct {
	IsStdin bool
	Stdin   io.Reader
}

type Terminal interface {
	ReadPassword(ctx context.Context) ([]byte, error)
}

type TokenManager interface {
	Set(token string) error
}
