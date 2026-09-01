package rmtoken

import "log/slog"

type Controller struct {
	param        *Param
	tokenManager TokenManager
}

func New(param *Param, tokenManager TokenManager) *Controller {
	return &Controller{
		param:        param,
		tokenManager: tokenManager,
	}
}

type Param struct{}

type TokenManager interface {
	Remove(logger *slog.Logger) error
}
