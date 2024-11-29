package vadu_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/contbank/vadu-sdk"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticationTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	ctx            context.Context
	session        *vadu.Session
	authentication *vadu.Authentication
	logger         *logrus.Logger
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func (s *AuthenticationTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	// Mock da configuração
	config := vadu.Config{
		ClientToken:   vadu.String("Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJWYWR1IiwidXNyIjoxMTM5NiwiZW1sIjoiY29uZmlnQGNvbnRiYW5rLmNvbS5iciIsImVtcCI6NjY1MTM2MDl9.U9DbXl6UbNPtv9_ZZjgdgodF-ISQIz_B1NPG0meXXXX"), // Substitua pelo token correto para testes reais
		Cookie:        vadu.String("mock-cookie-value"),
		LoginEndpoint: vadu.String("https://www.vadu.com.br/vadu.dll/Autenticacao/JSONPegarToken"),
	}
	// Configurar o logger
	s.logger = logrus.New()
	s.logger.SetFormatter(&logrus.JSONFormatter{}) // Formato estruturado JSON
	s.logger.SetOutput(ioutil.Discard)             // Descartar logs durante testes (opcional)

	session, err := vadu.NewSession(config)

	s.assert.NoError(err)

	// Cliente HTTP mockado
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.authentication = vadu.NewAuthentication(httpClient, *s.session, s.logger)
}

func (s *AuthenticationTestSuite) TestToken() {
	// Teste o método Token
	token, err := s.authentication.Token(s.ctx)
	s.assert.NotNil(token, "O token não pode ser nulo")
	s.assert.NoError(err, "Deveria obter o token sem erro")

}
func (s *AuthenticationTestSuite) TestInvalidToken() {
	// Configurar transporte mockado
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       ioutil.NopCloser(strings.NewReader(`{"erro":{"StatusCode":403,"Descricao":"Não autorizado!","Mensagem":"Token inválido"}}`)),
			}, nil
		},
	}

	// Criar cliente HTTP usando o transporte mockado
	mockClient := &http.Client{Transport: mockTransport}

	// Substituir o cliente HTTP na autenticação
	s.authentication = vadu.NewAuthentication(mockClient, *s.session, s.logger)

	s.session.ClientToken = "invalid-token"
	s.session.Cache.Delete("token") // Limpar cache

	token, err := s.authentication.Token(s.ctx)

	// Verificar o comportamento esperado
	s.Assert().Error(err, "Deveria retornar um erro para token inválido")
	s.Assert().Contains(err.Error(), "403", "Erro esperado deve indicar falha de autorização")
	s.Assert().Empty(token, "O token deveria estar vazio em caso de erro")
}
