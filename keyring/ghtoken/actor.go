package ghtoken

import (
	"context"
	"log/slog"
	"os"

	"github.com/suzuki-shunsuke/cobra-util/keyring/ghtoken/controller/rmtoken"
	"github.com/suzuki-shunsuke/cobra-util/keyring/ghtoken/controller/settoken"
)

type Actor struct {
	logger       *slog.Logger
	tokenService string
}

func NewActor(logger *slog.Logger, tokenService string) *Actor {
	return &Actor{
		logger:       logger,
		tokenService: tokenService,
	}
}

func (r *Actor) Set(ctx context.Context, input *InputSet) error {
	term := NewPasswordReader(os.Stdout)
	tokenManager := NewTokenManager(r.tokenService)
	ctrl := settoken.New(&settoken.Param{
		IsStdin: input.Stdin,
		Stdin:   os.Stdin,
	}, term, tokenManager)
	return ctrl.Set(ctx, r.logger) //nolint:wrapcheck
}

func (r *Actor) Remove() error {
	tokenManager := NewTokenManager(r.tokenService)
	ctrl := rmtoken.New(&rmtoken.Param{}, tokenManager)
	return ctrl.Remove(r.logger) //nolint:wrapcheck
}
