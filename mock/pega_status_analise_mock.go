package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// PegaStatusAnaliseMock retorna um mock para o método PegaStatusAnalise
func PegaStatusAnaliseMock() *http.Client {
	return &http.Client{
		Transport: &MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Simulando uma resposta bem-sucedida da API
				responseJSON := `{
					"quantidade_cnpj_cpf": 1,
					"quantidade_consultas_receita": 0,
					"percentual_consultas_receita": 0,
					"quantidade_cnpj_cpf_concluidos": 1,
					"percentual_concluido": 100,
					"finalizando_arquivo": false,
					"concluido": true
				}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(responseJSON))),
				}, nil
			},
		},
	}
}
