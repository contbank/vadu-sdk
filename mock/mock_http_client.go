package mock

import (
	"context"
	"log"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockAuthHTTPClient simula o cliente HTTP para as requisições de autenticação.
type MockAuthHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip implementa a interface http.RoundTripper.
// Ele é responsável por processar as requisições HTTP mockadas.
func (m *MockAuthHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	// Adicione um log para verificar se o mock está sendo chamado
	log.Printf("Mock RoundTrip chamado para URL: %s", req.URL.String())
	return m.DoFunc(req)
}

// MockAuthentication implementa AuthenticationInterface
type MockAuthentication struct {
	mock.Mock
}

// Token implementa o método da interface AuthenticationInterface
func (m *MockAuthentication) Token(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
