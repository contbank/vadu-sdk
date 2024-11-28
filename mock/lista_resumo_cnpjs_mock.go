package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// ListaResumoCNPJsMock retorna um mock para o método ListaResumoCNPJs
func ListaResumoCNPJsMock() *http.Client {
	return &http.Client{
		Transport: &MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Simulando uma resposta bem-sucedida da API com uma lista de resumos de CNPJs
				responseJSON := `[
					{
						"analise_id": 4768906,
						"analise_cnpj_cpf_id": 28883956,
						"cnpj_cpf": "98960887000164",
						"nome": "WEBSOLUTIONS LTDA",
						"erro": false,
						"alerta": true,
						"bloqueio": false,
						"rating": 500,
						"rating_sigla": "B (500)",
						"rating_descricao": "CCB--> 40% do faturamento mensal || ANT--> 20% do faturamento mensal",
						"rating2": 0,
						"rating2_sigla": "Fora de faixa: 0",
						"rating2_descricao": "",
						"nova_consulta_serasa": false,
						"nova_consulta_serasa_string_retorno": "",
						"flowSolicitacao_id": 0,
						"flowTarefaNome": "",
						"flowTarefa_id": 0,
						"origem_consulta_serasa": 0,
						"origem_consulta_serasa_texto": "Consulta Serasa"
					}
				]`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(responseJSON))),
				}, nil
			},
		},
	}
}
