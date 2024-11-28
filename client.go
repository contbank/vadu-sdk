package vadu

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// VaduClient estrutura principal para interagir com a API Vadu.
type VaduClient struct {
	httpClient *http.Client
	session    Session
}

// NewVaduClient cria uma nova instância do cliente da API Vadu.
func NewVaduClient(httpClient *http.Client, session Session) *VaduClient {
	return &VaduClient{
		httpClient: httpClient,
		session:    session,
	}
}

// ListaGruposAnalise lista os grupos de análise disponíveis na API do Vadu.
func (vc *VaduClient) ListaGruposAnalise(ctx context.Context) ([]GrupoAnalise, error) {
	url := fmt.Sprintf("%s/api-analise-bordero-config/v1/grupoanalise/cnpjcpf", vc.session.APIEndpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", vc.session.ClientToken))
	req.Header.Add("Cookie", vc.session.Cookie)

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro ao listar grupos de análise, status: %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	var grupos []GrupoAnalise
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &grupos); err != nil {
		return nil, err
	}

	return grupos, nil
}

// EnviaCNPJsParaAnalise envia uma lista de CNPJs para análise.
func (vc *VaduClient) EnviaCNPJsParaAnalise(cnpjEmpresa string, idGrupoAnalise int, listaCNPJCPF []string, postBack *PostBack) (*EnviaCNPJsResponse, error) {
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
		return nil, err
	}

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise", vc.session.APIEndpoint)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Definir os cabeçalhos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+vc.session.ClientToken)

	// Enviar a requisição usando o httpClient
	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Se o status da resposta não for 200, retornar erro
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("falha ao enviar CNPJs para análise, status: " + resp.Status)
	}

	// Decodificar a resposta
	var response EnviaCNPJsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// Retornar a resposta da API
	return &response, nil
}

// EnviaCNPJsComDadosParaAnalise envia uma lista de CNPJs com dados detalhados para análise.
func (vc *VaduClient) EnviaCNPJsComDadosParaAnalise(cnpjEmpresa string, idGrupoAnalise int, listaDados []DadosIntegracao) (*EnviaCNPJsResponse, error) {
	// Montar o corpo da requisição
	requestBody := EnviaCNPJsComDadosRequest{
		CNPJEmpresa:                 cnpjEmpresa,
		IDGrupoAnalise:              idGrupoAnalise,
		ListaCNPJCPFDadosIntegracao: listaDados,
	}

	// Converter o corpo para JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v2/erp/analise", vc.session.APIEndpoint)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Definir os cabeçalhos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+vc.session.ClientToken)

	// Enviar a requisição usando o httpClient
	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Se o status da resposta não for 200, retornar erro
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("falha ao enviar CNPJs com dados detalhados para análise, status: %d, resposta: %s", resp.StatusCode, string(respBody))
	}

	// Decodificar a resposta
	var response EnviaCNPJsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// Retornar a resposta da API
	return &response, nil
}

// PegaStatusAnalise busca o status de uma análise pelo ID fornecido
func (vc *VaduClient) PegaStatusAnalise(idAnalise int) (*StatusAnalise, error) {
	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/status/analise/id/%d", vc.session.APIEndpoint, idAnalise)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Definir os cabeçalhos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+vc.session.ClientToken)

	// Realiza a requisição
	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requisição falhou com status %d", resp.StatusCode)
	}

	// Decodifica a resposta
	var status StatusAnalise
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &status, nil
}

// PegaResumoAnalise busca o resumo de uma análise pelo ID fornecido
func (vc *VaduClient) PegaResumoAnalise(idAnalise int) (*ResumoAnalise, error) {
	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise/id/%d", vc.session.APIEndpoint, idAnalise)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Adiciona o token de autorização ao cabeçalho
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", vc.session.ClientToken))

	// Realiza a requisição
	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requisição falhou com status %d", resp.StatusCode)
	}

	// Decodifica a resposta
	var resumo ResumoAnalise
	err = json.NewDecoder(resp.Body).Decode(&resumo)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &resumo, nil
}

// ListaResumoCNPJs busca os resumos dos CNPJs analisados para uma análise pelo ID fornecido
func (vc *VaduClient) ListaResumoCNPJs(idAnalise int) ([]ResumoCNPJ, error) {

	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise/id/%d/cnpjcpf", vc.session.APIEndpoint, idAnalise)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Adiciona o token de autorização ao cabeçalho
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", vc.session.ClientToken))

	// Realiza a requisição
	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requisição falhou com status %d", resp.StatusCode)
	}

	// Decodifica a resposta
	var resumos []ResumoCNPJ
	err = json.NewDecoder(resp.Body).Decode(&resumos)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return resumos, nil
}

func (vc *VaduClient) ListaResumoCNPJsDetalhado(analiseID int) ([]ResumoCNPJDatalhado, error) {
	url := fmt.Sprintf("%s/api-analise-cnpjcpf/v1/erp/analise/id/%d/cnpjcpf/detalhado", vc.session.APIEndpoint, analiseID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", vc.session.ClientToken))

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var resumos []ResumoCNPJDatalhado
	if err := json.NewDecoder(resp.Body).Decode(&resumos); err != nil {
		return nil, err
	}

	return resumos, nil
}
