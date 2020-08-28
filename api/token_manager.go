package api

import (
	"sync"
	"time"
)

var tokenManager = &TokenManager{
	tokenSet: map[string]*TenantAccessToken{},
}

type TokenManager struct {
	mu       sync.Mutex
	tokenSet map[string]*TenantAccessToken
}

func (t *TokenManager) Get(key string) *TenantAccessToken {
	t.mu.Lock()
	defer t.mu.Unlock()

	token := t.tokenSet[key]
	if token != nil {
		usable := time.Now().Before(token.ExpiredAt)
		if usable {
			return token
		} else {
			t.del(key)
		}
	}
	return nil
}

func (t *TokenManager) Set(key string, token *TenantAccessToken) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tokenSet[key] = token
}

func (t *TokenManager) del(key string) {
	delete(t.tokenSet, key)
}
