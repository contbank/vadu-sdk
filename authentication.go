package vadu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Authentication define a estrutura para autenticação no Vadu SDK.
type Authentication struct {
	session    Session
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewAuthentication inicializa uma nova instância de Authentication.
func NewAuthentication(httpClient *http.Client, session Session, logger *logrus.Logger) *Authentication {
	return &Authentication{
		session:    session,
		httpClient: httpClient,
		logger:     logger,
	}
}

// AuthenticationResponse representa a resposta da API de autenticação.
type AuthenticationResponse struct {
	Token string `json:"token"`
}

// Error padrão para login.
var ErrDefaultLogin = errors.New("falha ao autenticar")

// login realiza o request à API do Vadu para obter o token de autenticação.
func (a *Authentication) login(ctx context.Context) (*AuthenticationResponse, error) {
	a.logger.WithFields(logrus.Fields{
		"endpoint": a.session.LoginEndpoint,
	}).Info("Iniciando login no Vadu")

	if a.session.ClientToken == "" {
		err := errors.New("ClientToken não fornecido")
		a.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Erro ao autenticar")
		return nil, err
	}

	// Cria uma requisição HTTP com timeout de contexto.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", a.session.LoginEndpoint, nil)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Erro ao criar requisição HTTP")
		return nil, err
	}

	// Adiciona os headers necessários.
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s", a.session.ClientToken))
	req.Header.Add("Cookie", a.session.Cookie)

	// Envia o request.
	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"endpoint": a.session.LoginEndpoint,
			"error":    err,
		}).Error("Erro ao enviar requisição HTTP")
		return nil, err
	}
	defer resp.Body.Close()

	// Verifica o status da resposta.
	if resp.StatusCode == http.StatusUnauthorized {
		a.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"endpoint":    a.session.LoginEndpoint,
		}).Warn("Autenticação não autorizada (401)")
		return nil, fmt.Errorf("não autorizado (401): %w", ErrDefaultLogin)
	} else if resp.StatusCode >= 500 {
		a.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"endpoint":    a.session.LoginEndpoint,
		}).Error("Erro no servidor do Vadu")
		return nil, fmt.Errorf("erro no servidor (500): %s", resp.Status)
	} else if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		a.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(respBody),
			"endpoint":    a.session.LoginEndpoint,
		}).Error("Erro ao autenticar")
		return nil, fmt.Errorf("erro ao autenticar, status: %d", resp.StatusCode)
	}

	// Lê e parseia a resposta.
	var response AuthenticationResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		a.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Erro ao parsear resposta JSON")
		return nil, err
	}

	a.logger.WithFields(logrus.Fields{
		"endpoint": a.session.LoginEndpoint,
	}).Info("Login realizado com sucesso")
	return &response, nil
}

// Token retorna o token de autenticação armazenado no cache ou faz login para obtê-lo.
func (a *Authentication) Token(ctx context.Context) (string, error) {
	// Verifica se o token já está em cache.
	if token, found := a.session.Cache.Get("token"); found {
		a.logger.WithFields(logrus.Fields{
			"cache": "hit",
		}).Info("Token obtido do cache")
		return token.(string), nil
	}

	// Realiza o login para obter um novo token.
	response, err := a.login(ctx)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Erro ao obter token via login")
		return "", err
	}

	// Define a validade do token como 24 horas.
	const tokenValidity = 24 * time.Hour
	a.session.Cache.Set("token", response.Token, tokenValidity)

	// Masca o token para segurança ao logar.
	maskedToken := response.Token[:5] + "..."
	a.logger.WithFields(logrus.Fields{
		"cache": "set",
		"token": maskedToken,
	}).Info("Novo token autenticado e armazenado no cache")

	return response.Token, nil
}
