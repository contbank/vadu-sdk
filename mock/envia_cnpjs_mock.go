package mock

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Mock para a função EnviaCNPJsParaAnalise
func EnviaCNPJsParaAnaliseMock() *http.Client {
	return &http.Client{
		Transport: &MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(strings.NewReader(`
					{
						"analise_id": 4768906,
						"quantidade_cnpj": 1,
						"quantidade_cpf": 0,
						"usuario": "Contbank - Usuário p/Integração Não excluir",
						"data_hora_envio": "2024-11-25T18:28:44.523Z",
						"id_grupo_analise": 10802,
						"nome_lote": "",
						"nome_grupo_analise": "Análises CNPJs"
					}`)),
				}, nil
			},
		},
	}
}
