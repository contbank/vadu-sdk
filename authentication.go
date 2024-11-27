package vadu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Authentication define a estrutura para autenticação no Vadu SDK.
type Authentication struct {
	session    Session
	httpClient *http.Client
}

// NewAuthentication inicializa uma nova instância de Authentication.
func NewAuthentication(httpClient *http.Client, session Session) *Authentication {
	return &Authentication{
		session:    session,
		httpClient: httpClient,
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
	if a.session.ClientToken == "" {
		return nil, errors.New("ClientToken não fornecido")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", a.session.LoginEndpoint, nil)
	if err != nil {
		return nil, err
	}

	// Adiciona os headers necessários
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s", a.session.ClientToken))
	req.Header.Add("Cookie", a.session.Cookie)

	// Envia o request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao autenticar, status: %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Lê e parseia a resposta
	var response AuthenticationResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Token retorna o token de autenticação armazenado no cache ou faz login para obtê-lo.
func (a *Authentication) Token(ctx context.Context) (string, error) {
	if token, found := a.session.Cache.Get("token"); found {
		return token.(string), nil
	}

	response, err := a.login(ctx)
	if err != nil {
		return "", err
	}

	// Armazena o token no cache, com duração baseada no tempo de expiração padrão
	a.session.Cache.Set("token", response.Token, time.Minute*30)

	return response.Token, nil
}
