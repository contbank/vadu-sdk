package vadu

import (
	"net/http"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
)

// Config contém as configurações necessárias para inicializar uma sessão.
type Config struct {
	APIEndpoint   *string        // URL do API
	LoginEndpoint *string        // URL de autenticação
	ClientToken   *string        // Token do cliente
	Cookie        *string        // Cookie de autenticação
	Cache         *cache.Cache   // Cache para armazenar o token
	HTTPClient    *http.Client   // Cliente HTTP personalizado
	TokenTTL      *time.Duration // Tempo de expiração do token (opcional)
}

// Session representa a sessão autenticada com as configurações da API do Vadu.
type Session struct {
	APIEndpoint   string        // URL do API
	LoginEndpoint string        // URL para autenticação
	ClientToken   string        // Token do cliente
	Cookie        string        // Cookie de autenticação
	Cache         *cache.Cache  // Cache para tokens
	HTTPClient    *http.Client  // Cliente HTTP
	TokenTTL      time.Duration // Tempo de expiração do token
}

// NewSession cria uma nova instância de `Session` com base nas configurações fornecidas.
func NewSession(config Config) (*Session, error) {
	if config.APIEndpoint == nil {
		config.APIEndpoint = String("https://www.vadu.com.br")
	}

	if config.LoginEndpoint == nil {
		config.LoginEndpoint = String("https://www.vadu.com.br/vadu.dll/Autenticacao/JSONPegarToken")
	}

	if config.ClientToken == nil {
		config.ClientToken = String(os.Getenv("VADU_CLIENT_TOKEN"))
	}

	if config.Cookie == nil {
		config.Cookie = String(os.Getenv("VADU_COOKIE"))
	}

	if config.Cache == nil {
		config.Cache = cache.New(10*time.Minute, 1*time.Minute)
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}

	if config.TokenTTL == nil {
		defaultTTL := 30 * time.Minute
		config.TokenTTL = &defaultTTL
	}

	// Inicializa a sessão
	return &Session{
		LoginEndpoint: *config.LoginEndpoint,
		ClientToken:   *config.ClientToken,
		Cookie:        *config.Cookie,
		Cache:         config.Cache,
		HTTPClient:    config.HTTPClient,
		TokenTTL:      *config.TokenTTL,
	}, nil
}
