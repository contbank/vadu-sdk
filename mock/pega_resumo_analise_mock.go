package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// PegaResumoAnaliseMock retorna um mock para o método PegaResumoAnalise
func PegaResumoAnaliseMock() *http.Client {
	return &http.Client{
		Transport: &MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Simulando uma resposta bem-sucedida da API
				responseJSON := `{
					"analise_id": 4768906,
					"quantidade_cnpj": 1,
					"quantidade_cpf": 0,
					"cnpj_empresa": "33011770000199",
					"usuario": "Contbank - Usuário p/Integração Não excluir",
					"data_hora_envio": "2024-11-25T18:28:44.523Z",
					"data_hora_conclusao": "2024-11-25T18:28:50.483Z",
					"concluido": true,
					"erro": false,
					"alerta": true,
					"bloqueio": false,
					"quantidade_cnpj_alerta": 1,
					"quantidade_cnpj_bloqueio": 0,
					"quantidade_cpf_alerta": 0,
					"quantidade_cpf_bloqueio": 0,
					"id_grupo_analise": 10802,
					"nome_grupo_analise": "Análises CNPJs",
					"rating_valor": 500,
					"rating_sigla": "B (500)",
					"rating_descricao": "CCB--> 40% do faturamento mensal || ANT--> 20% do faturamento mensal",
					"rating2_valor": 0,
					"rating2_sigla": "Fora de faixa: 0",
					"rating2_descricao": "Rating fora de faixa",
					"nome_lote": ""
				}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(responseJSON))),
				}, nil
			},
		},
	}
}
