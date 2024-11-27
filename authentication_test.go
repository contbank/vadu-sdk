package vadu_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/vadu-sdk/vadu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticationTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	ctx            context.Context
	session        *vadu.Session
	authentication *vadu.Authentication
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}

func (s *AuthenticationTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	// Mock da configuração
	config := vadu.Config{
		ClientToken:   vadu.String("mock-client-token"), // Substitua pelo token correto para testes reais
		Cookie:        vadu.String("mock-cookie-value"),
		LoginEndpoint: vadu.String("https://www.vadu.com.br/vadu.dll/Autenticacao/JSONPegarToken"),
	}

	session, err := vadu.NewSession(config)

	s.assert.NoError(err)

	// Cliente HTTP mockado
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.authentication = vadu.NewAuthentication(httpClient, *s.session)
}

func (s *AuthenticationTestSuite) TestToken() {
	// Teste o método Token
	token, err := s.authentication.Token(s.ctx)

	s.assert.NoError(err, "Deveria obter o token sem erro")
	s.assert.Contains(token, "Bearer", "O token deveria conter a palavra 'Bearer'")
}

func (s *AuthenticationTestSuite) TestInvalidToken() {
	// Alterar o token para um valor inválido no contexto de teste
	s.session.ClientToken = "invalid-token"

	// O método Token deve retornar um erro
	token, err := s.authentication.Token(s.ctx)

	s.assert.Error(err, "Deveria retornar um erro para token inválido")
	s.assert.Empty(token, "O token deveria estar vazio em caso de erro")
}
