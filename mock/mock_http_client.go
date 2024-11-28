package mock

import (
	"log"
	"net/http"
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
