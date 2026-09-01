package ghtoken

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	keyName = "GITHUB_TOKEN"
)

type TokenManager struct {
	keyService string
}

func NewTokenManager(keyService string) *TokenManager {
	return &TokenManager{
		keyService: keyService,
	}
}

func (tm *TokenManager) Set(token string) error {
	if err := keyring.Set(tm.keyService, keyName, token); err != nil {
		return fmt.Errorf("set a GitHub Access token in keyring: %w", err)
	}
	return nil
}

func (tm *TokenManager) Remove(logger *slog.Logger) error {
	if err := keyring.Delete(tm.keyService, keyName); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			slogerr.WithError(logger, err).Warn("remove a GitHub Access token from keyring")
			return nil
		}
		return fmt.Errorf("delete a GitHub Access token from keyring: %w", err)
	}
	return nil
}

type TokenSource struct {
	token      *oauth2.Token
	logger     *slog.Logger
	mutex      *sync.RWMutex
	keyService string
}

func NewTokenSource(logger *slog.Logger, keyService string) *TokenSource {
	return &TokenSource{
		logger:     logger,
		mutex:      &sync.RWMutex{},
		keyService: keyService,
	}
}

func (ks *TokenSource) Token() (*oauth2.Token, error) {
	ks.mutex.RLock()
	token := ks.token
	ks.mutex.RUnlock()
	if token != nil {
		return token, nil
	}
	ks.logger.Debug("getting a GitHub Access token from keyring")
	s, err := keyring.Get(ks.keyService, keyName)
	if err != nil {
		return nil, fmt.Errorf("get a GitHub Access token from keyring: %w", err)
	}
	ks.logger.Debug("got a GitHub Access token from keyring")
	token = &oauth2.Token{
		AccessToken: s,
	}
	ks.mutex.Lock()
	ks.token = token
	ks.mutex.Unlock()
	return token, nil
}
