package mock

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Mock para a função ListaGruposAnalise
func ListaGruposAnaliseMock() *http.Client {
	return &http.Client{
		Transport: &MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(strings.NewReader(`
					[
						{
							"id_grupo_analise": 10802,
							"nome_grupo_analise": "Análises CNPJs",
							"rating_start": 1000,
							"rating_minimo": -15650,
							"rating_maximo": 0,
							"quantidade_analises": 1,
							"quantidade_regras": 28,
							"quantidade_validacoes": 79
						}
					]`)),
				}, nil
			},
		},
	}
}
