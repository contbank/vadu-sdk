package vadu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// VaduClient estrutura principal para interagir com a API Vadu.
type VaduClient struct {
	httpClient *http.Client
	session    Session
	logger     *logrus.Logger
}

// NewVaduClient cria uma nova instância do cliente da API Vadu.
func NewVaduClient(httpClient *http.Client, session Session, logger *logrus.Logger) *VaduClient {
	return &VaduClient{
		httpClient: httpClient,
		session:    session,
		logger:     logger,
	}
}

// ListaGruposAnalise lista os grupos de análise disponíveis na API do Vadu.
func (vc *VaduClient) ListaGruposAnalise(ctx context.Context, auth AuthenticationInterface) ([]GrupoAnalise, error) {

	// Obtenha o token dinamicamente
	token, err := auth.Token(ctx)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao obter token de autenticação")
		return nil, fmt.Errorf("falha ao autenticar: %w", err)
	}

	// Obtenha o token dinamicamente
	url := fmt.Sprintf("%s/api-analise-bordero-config/v1/grupoanalise/cnpjcpf", vc.session.APIEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao criar requisição para listar grupos de análise")
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Cookie", vc.session.Cookie)

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		vc.logger.WithError(err).Error("Erro de conexão ao tentar listar grupos de análise")
		return nil, fmt.Errorf("erro de conexão com o servidor: %w", err)
	}
	defer resp.Body.Close()

	// Logar o status HTTP
	vc.logger.WithFields(logrus.Fields{
		"url":        url,
		"statusCode": resp.StatusCode,
	}).Info("Resposta da API recebida")

	// Verificar status HTTP
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		vc.logger.WithFields(logrus.Fields{
			"url":        url,
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Erro ao listar grupos de análise")
		return nil, fmt.Errorf("erro ao listar grupos de análise: status %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Ler corpo da resposta
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, fmt.Errorf("erro ao ler resposta da API: %w", err)
	}

	// Verificar formato JSON
	var grupos []GrupoAnalise
	if err := json.Unmarshal(respBody, &grupos); err != nil {
		vc.logger.WithError(err).Error("Erro ao decodificar JSON da resposta")
		return nil, fmt.Errorf("erro no formato da resposta da API: %w", err)
	}

	// Sucesso
	vc.logger.WithFields(logrus.Fields{
		"url":    url,
		"grupos": len(grupos),
	}).Info("Grupos de análise listados com sucesso")
	return grupos, nil
}

// EnviaCNPJsParaAnalise envia uma lista de CNPJs para análise com validações e logs.
func (vc *VaduClient) EnviaCNPJsParaAnalise(ctx context.Context, cnpjEmpresa string, idGrupoAnalise int, listaCNPJCPF []string, postBack *PostBack, auth AuthenticationInterface) (*EnviaCNPJsResponse, error) {

	// Obtenha o token dinamicamente
	token, err := auth.Token(ctx)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao obter token de autenticação")
		return nil, fmt.Errorf("falha ao autenticar: %w", err)
	}
	// Validar o número de CNPJs
	if len(listaCNPJCPF) > 2000 {
		vc.logger.WithFields(logrus.Fields{
			"cnpjEmpresa":     cnpjEmpresa,
			"idGrupoAnalise":  idGrupoAnalise,
			"quantidadeCNPJs": len(listaCNPJCPF),
		}).Error("Número máximo de CNPJs excedido")
		return nil, fmt.Errorf("não é permitido enviar mais de 2000 CNPJs por requisição")
	}

	// Montar o corpo da requisição
	requestBody := EnviaCNPJsRequest{
		CNPJEmpresa:    cnpjEmpresa,
		IDGrupoAnalise: idGrupoAnalise,
		ListaCNPJCPF:   listaCNPJCPF,
		PostBack:       postBack, // postBack pode ser nil
	}

	// Converter o corpo para JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao converter o corpo da requisição para JSON")
		return nil, fmt.Errorf("erro ao preparar o payload: %w", err)
	}

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise", vc.session.APIEndpoint)

	// Logar o payload enviado
	vc.logger.WithFields(logrus.Fields{
		"url":     url,
		"payload": string(jsonData),
	}).Info("Enviando requisição para análise de CNPJs")

	var resp *http.Response
	maxRetries := 3 // Número máximo de tentativas
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Criar a requisição HTTP
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			vc.logger.WithError(err).Error("Erro ao criar requisição HTTP")
			return nil, fmt.Errorf("erro ao criar a requisição: %w", err)
		}

		// Definir os cabeçalhos
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		// Enviar a requisição usando o httpClient
		resp, err = vc.httpClient.Do(req)
		if err != nil {
			vc.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Error("Erro ao enviar a requisição")
			if attempt == maxRetries {
				return nil, fmt.Errorf("falha ao conectar com o servidor após %d tentativas: %w", maxRetries, err)
			}
			continue // Tentar novamente
		}
		break
	}
	defer resp.Body.Close()

	// Logar o status HTTP da resposta
	vc.logger.WithFields(logrus.Fields{
		"statusCode": resp.StatusCode,
		"url":        url,
	}).Info("Resposta recebida da API")

	// Se o status da resposta não for 200, retornar erro
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		vc.logger.WithFields(logrus.Fields{
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Falha ao enviar CNPJs para análise")
		return nil, fmt.Errorf("erro ao enviar CNPJs para análise: status %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Decodificar a resposta
	var response EnviaCNPJsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao decodificar JSON da resposta")
		return nil, fmt.Errorf("erro no formato da resposta da API: %w", err)
	}

	// Logar a resposta bem-sucedida
	vc.logger.WithFields(logrus.Fields{
		"statusCode": resp.StatusCode,
		"response":   response,
	}).Info("CNPJs enviados para análise com sucesso")

	// Retornar a resposta da API
	return &response, nil
}

// EnviaCNPJsComDadosParaAnalise envia uma lista de CNPJs com dados detalhados para análise com validações e logs.
func (vc *VaduClient) EnviaCNPJsComDadosParaAnalise(ctx context.Context, cnpjEmpresa string, idGrupoAnalise int, listaDados []DadosIntegracao, auth AuthenticationInterface) (*EnviaCNPJsResponse, error) {
	// Obtenha o token dinamicamente
	token, err := auth.Token(ctx)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao obter token de autenticação")
		return nil, fmt.Errorf("falha ao autenticar: %w", err)
	}
	// Validar o número de CNPJs
	if len(listaDados) > 100 {
		vc.logger.WithFields(logrus.Fields{
			"cnpjEmpresa":     cnpjEmpresa,
			"idGrupoAnalise":  idGrupoAnalise,
			"quantidadeCNPJs": len(listaDados),
		}).Error("Número máximo de CNPJs excedido")
		return nil, fmt.Errorf("não é permitido enviar mais de 100 CNPJs por requisição")
	}

	// Montar o corpo da requisição
	requestBody := EnviaCNPJsComDadosRequest{
		CNPJEmpresa:                 cnpjEmpresa,
		IDGrupoAnalise:              idGrupoAnalise,
		ListaCNPJCPFDadosIntegracao: listaDados,
	}

	// Converter o corpo para JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao converter o corpo da requisição para JSON")
		return nil, fmt.Errorf("erro ao preparar o payload: %w", err)
	}

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v2/erp/analise", vc.session.APIEndpoint)

	// Logar o payload enviado
	vc.logger.WithFields(logrus.Fields{
		"url":     url,
		"payload": string(jsonData),
	}).Info("Enviando requisição para análise detalhada de CNPJs")

	var resp *http.Response
	maxRetries := 3 // Número máximo de tentativas
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Criar a requisição HTTP
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			vc.logger.WithError(err).Error("Erro ao criar requisição HTTP")
			return nil, fmt.Errorf("erro ao criar a requisição: %w", err)
		}

		// Definir os cabeçalhos
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		// Enviar a requisição usando o httpClient
		resp, err = vc.httpClient.Do(req)
		if err != nil {
			vc.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Error("Erro ao enviar a requisição")
			if attempt == maxRetries {
				return nil, fmt.Errorf("falha ao conectar com o servidor após %d tentativas: %w", maxRetries, err)
			}
			continue // Tentar novamente
		}
		break
	}
	defer resp.Body.Close()

	// Logar o status HTTP da resposta
	vc.logger.WithFields(logrus.Fields{
		"statusCode": resp.StatusCode,
		"url":        url,
	}).Info("Resposta recebida da API")

	// Se o status da resposta não for 200, retornar erro
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		vc.logger.WithFields(logrus.Fields{
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Falha ao enviar CNPJs com dados detalhados para análise")
		return nil, fmt.Errorf("erro ao enviar CNPJs: status %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Decodificar a resposta
	var response EnviaCNPJsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao decodificar JSON da resposta")
		return nil, fmt.Errorf("erro no formato da resposta da API: %w", err)
	}

	// Logar a resposta bem-sucedida
	vc.logger.WithFields(logrus.Fields{
		"statusCode": resp.StatusCode,
		"response":   response,
	}).Info("CNPJs enviados para análise detalhada com sucesso")

	// Retornar a resposta da API
	return &response, nil
}

// PegaStatusAnalise busca o status de uma análise pelo ID fornecido, com validações, logs e timeout configurado.
func (vc *VaduClient) PegaStatusAnalise(ctx context.Context, analiseID int, auth AuthenticationInterface) (*StatusAnalise, error) {
	// Obtenha o token dinamicamente
	token, err := auth.Token(ctx)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao obter token de autenticação")
		return nil, fmt.Errorf("falha ao autenticar: %w", err)
	}
	// Validar se o ID da análise é positivo e numérico
	if analiseID <= 0 {
		vc.logger.WithFields(logrus.Fields{
			"analiseID": analiseID,
		}).Error("ID de análise inválido")
		return nil, fmt.Errorf("analiseID deve ser um número positivo")
	}

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/status/analise/id/%d", vc.session.APIEndpoint, analiseID)

	// Logar o ID da análise sendo consultado
	vc.logger.WithFields(logrus.Fields{
		"analiseID": analiseID,
		"url":       url,
	}).Info("Consultando status da análise")

	// Configurar contexto e timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Criar a requisição HTTP com o contexto configurado
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao criar requisição HTTP")
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Definir os cabeçalhos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Enviar a requisição utilizando o httpClient configurado na struct
	var resp *http.Response
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = vc.httpClient.Do(req) // Usando httpClient da struct
		if err != nil {
			vc.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Error("Erro ao realizar requisição")
			if attempt == 3 {
				return nil, fmt.Errorf("falha ao conectar ao servidor após 3 tentativas: %w", err)
			}
			continue // Tentar novamente
		}
		break
	}
	defer resp.Body.Close()

	// Logar o status da resposta
	vc.logger.WithFields(logrus.Fields{
		"analiseID":  analiseID,
		"statusCode": resp.StatusCode,
	}).Info("Resposta recebida da API")

	// Se o status da resposta não for 200, retornar erro
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		vc.logger.WithFields(logrus.Fields{
			"analiseID":  analiseID,
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Falha ao consultar status da análise")
		return nil, fmt.Errorf("erro ao consultar status da análise: status %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Decodificar a resposta
	var status StatusAnalise
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao decodificar resposta da API")
		return nil, fmt.Errorf("erro no formato da resposta da API: %w", err)
	}

	// Logar a resposta bem-sucedida
	vc.logger.WithFields(logrus.Fields{
		"analiseID": analiseID,
		"status":    status,
	}).Info("Status da análise consultado com sucesso")

	// Retornar o status da análise
	return &status, nil
}

// PegaResumoAnalise busca o resumo de uma análise pelo ID fornecido, com validações, logs e timeout configurado.
func (vc *VaduClient) PegaResumoAnalise(ctx context.Context, analiseID int, auth AuthenticationInterface) (*ResumoAnalise, error) {
	// Obtenha o token dinamicamente
	token, err := auth.Token(ctx)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao obter token de autenticação")
		return nil, fmt.Errorf("falha ao autenticar: %w", err)
	}
	// Validação do ID de Análise - Certifica-se de que o ID é válido e numérico
	if analiseID <= 0 {
		vc.logger.WithFields(logrus.Fields{
			"analiseID": analiseID,
		}).Error("ID de análise inválido")
		return nil, fmt.Errorf("analiseID deve ser um número positivo")
	}

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise/id/%d", vc.session.APIEndpoint, analiseID)

	// Logar o ID da análise sendo consultado
	vc.logger.WithFields(logrus.Fields{
		"analiseID": analiseID,
		"url":       url,
	}).Info("Consultando resumo da análise")

	// Configurar contexto e timeout para a requisição
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Criar a requisição HTTP com o contexto configurado
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao criar requisição HTTP")
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Definir os cabeçalhos
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Enviar a requisição utilizando o httpClient da struct
	var resp *http.Response
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = vc.httpClient.Do(req) // Usando httpClient da struct
		if err != nil {
			vc.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Error("Erro ao realizar requisição")
			if attempt == 3 {
				return nil, fmt.Errorf("falha ao conectar ao servidor após 3 tentativas: %w", err)
			}
			continue // Tentar novamente
		}
		break
	}
	defer resp.Body.Close()

	// Logar o status da resposta
	vc.logger.WithFields(logrus.Fields{
		"analiseID":  analiseID,
		"statusCode": resp.StatusCode,
	}).Info("Resposta recebida da API")

	// Se o status da resposta não for 200, retornar erro
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		vc.logger.WithFields(logrus.Fields{
			"analiseID":  analiseID,
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Falha ao consultar resumo da análise")
		return nil, fmt.Errorf("erro ao consultar resumo da análise: status %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Decodificar a resposta
	var resumo ResumoAnalise
	err = json.NewDecoder(resp.Body).Decode(&resumo)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao decodificar resposta da API")
		return nil, fmt.Errorf("erro no formato da resposta da API: %w", err)
	}

	// Logar a resposta bem-sucedida
	vc.logger.WithFields(logrus.Fields{
		"analiseID": analiseID,
		"resumo":    resumo,
	}).Info("Resumo da análise consultado com sucesso")

	// Retornar o resumo da análise
	return &resumo, nil
}

// / ListaResumoCNPJs busca os resumos dos CNPJs analisados para uma análise pelo ID fornecido, com validações, logs e timeout configurado.
func (vc *VaduClient) ListaResumoCNPJs(ctx context.Context, analiseID int, auth AuthenticationInterface) ([]ResumoCNPJ, error) {
	// Obtenha o token dinamicamente
	token, err := auth.Token(ctx)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao obter token de autenticação")
		return nil, fmt.Errorf("falha ao autenticar: %w", err)
	}
	// Validação do ID de Análise - Certifica-se de que o ID é válido e numérico
	if analiseID <= 0 {
		vc.logger.WithFields(logrus.Fields{
			"analiseID": analiseID,
		}).Error("ID de análise inválido")
		return nil, fmt.Errorf("analiseID deve ser um número positivo")
	}

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise/id/%d/cnpjcpf", vc.session.APIEndpoint, analiseID)

	// Logar o ID da análise sendo consultado
	vc.logger.WithFields(logrus.Fields{
		"analiseID": analiseID,
		"url":       url,
	}).Info("Consultando resumo dos CNPJs para análise")

	// Configurar contexto e timeout para a requisição
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Criar a requisição HTTP com o contexto configurado
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao criar requisição HTTP")
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Definir os cabeçalhos
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Enviar a requisição utilizando o httpClient da struct
	var resp *http.Response
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err = vc.httpClient.Do(req) // Usando httpClient da struct
		if err != nil {
			vc.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Error("Erro ao realizar requisição")
			if attempt == 3 {
				return nil, fmt.Errorf("falha ao conectar ao servidor após 3 tentativas: %w", err)
			}
			continue // Tentar novamente
		}
		break
	}
	defer resp.Body.Close()

	// Logar o status da resposta
	vc.logger.WithFields(logrus.Fields{
		"analiseID":  analiseID,
		"statusCode": resp.StatusCode,
	}).Info("Resposta recebida da API")

	// Se o status da resposta não for 200, retornar erro
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		vc.logger.WithFields(logrus.Fields{
			"analiseID":  analiseID,
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Falha ao consultar resumo dos CNPJs")
		return nil, fmt.Errorf("erro ao consultar resumo dos CNPJs: status %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Decodificar a resposta
	var resumos []ResumoCNPJ
	err = json.NewDecoder(resp.Body).Decode(&resumos)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao decodificar resposta da API")
		return nil, fmt.Errorf("erro no formato da resposta da API: %w", err)
	}

	// Logar a resposta bem-sucedida
	vc.logger.WithFields(logrus.Fields{
		"analiseID": analiseID,
		"resumos":   resumos,
	}).Info("Resumos dos CNPJs consultados com sucesso")

	// Retornar os resumos dos CNPJs
	return resumos, nil
}

// ListaResumoCNPJsDetalhado busca os resumos detalhados dos CNPJs analisados para uma análise pelo ID fornecido,
// com validação de entrada, resiliência, e logs detalhados.

func (vc *VaduClient) ListaResumoCNPJsDetalhado(ctx context.Context, analiseID int, auth AuthenticationInterface) ([]ResumoCNPJDatalhado, error) {
	// Obtenha o token dinamicamente
	token, err := auth.Token(ctx)
	if err != nil {
		vc.logger.WithError(err).Error("Erro ao obter token de autenticação")
		return nil, fmt.Errorf("falha ao autenticar: %w", err)
	}
	// Validação de Entrada: Verifica se o analiseID é um número válido e obrigatório
	if analiseID <= 0 {
		vc.logger.WithFields(logrus.Fields{
			"analiseID": analiseID,
		}).Error("ID de análise inválido")
		return nil, fmt.Errorf("analiseID deve ser um número positivo")
	}

	// Formatar a URL para a requisição
	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise/id/%d/cnpjcpf/detalhado", vc.session.APIEndpoint, analiseID)

	// Log de início da consulta
	vc.logger.WithFields(logrus.Fields{
		"analiseID": analiseID,
		"url":       url,
	}).Info("Consultando resumo detalhado dos CNPJs")

	// Configuração de timeout e cliente HTTP com retries
	maxRetries := 5
	retryDelay := 500 * time.Millisecond

	var resp *http.Response

	// Tentativas com backoff exponencial
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Criar a requisição HTTP
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			vc.logger.WithError(err).Error("Erro ao criar requisição HTTP")
			return nil, fmt.Errorf("erro ao criar requisição: %w", err)
		}

		// Definir os cabeçalhos
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// Enviar a requisição usando o httpClient da struct
		resp, err = vc.httpClient.Do(req) // Usando httpClient da struct
		if err != nil {
			vc.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Error("Erro ao realizar requisição")
			if attempt == maxRetries {
				return nil, fmt.Errorf("falha ao conectar ao servidor após %d tentativas: %w", maxRetries, err)
			}
			// Aguarda antes de tentar novamente com backoff exponencial
			time.Sleep(retryDelay)
			retryDelay *= 2 // Backoff exponencial
			continue
		}
		break
	}
	defer resp.Body.Close()

	// Log do tempo de resposta da API
	vc.logger.WithFields(logrus.Fields{
		"analiseID":  analiseID,
		"statusCode": resp.StatusCode,
	}).Infof("Resposta recebida da API - Tempo de resposta: %v", resp.Status)

	// Se o status da resposta não for OK (200), loga o erro e retorna
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		vc.logger.WithFields(logrus.Fields{
			"analiseID":  analiseID,
			"statusCode": resp.StatusCode,
			"response":   string(respBody),
		}).Error("Falha ao consultar resumo detalhado dos CNPJs")
		return nil, fmt.Errorf("erro ao consultar resumo detalhado dos CNPJs: status %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Decodificar a resposta da API
	var resumos []ResumoCNPJDatalhado
	if err := json.NewDecoder(resp.Body).Decode(&resumos); err != nil {
		vc.logger.WithError(err).Error("Erro ao decodificar resposta da API")
		return nil, fmt.Errorf("erro ao decodificar a resposta da API: %w", err)
	}

	// Filtrar os resumos para retornar apenas logs com erro ou alerta
	var filteredResumos []ResumoCNPJDatalhado
	for _, resumo := range resumos {
		var filteredLogs []LogAnalise
		for _, log := range resumo.Logs {
			if log.Erro || log.Alerta {
				filteredLogs = append(filteredLogs, log)
			}
		}

		// Se houver logs filtrados, adicionar o resumo à lista
		if len(filteredLogs) > 0 {
			resumo.Logs = filteredLogs
			filteredResumos = append(filteredResumos, resumo)
		}
	}

	// Log de sucesso com os dados filtrados
	vc.logger.WithFields(logrus.Fields{
		"analiseID":  analiseID,
		"resumos":    filteredResumos,
		"logs_count": len(filteredResumos),
	}).Info("Consulta dos resumos detalhados bem-sucedida")

	// Retornar os resumos detalhados filtrados
	return filteredResumos, nil
}
