package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// ListaResumoCNPJsDetalhadoMock retorna um mock para a API detalhada
func ListaResumoCNPJsDetalhadoMock() *http.Client {
	return &http.Client{
		Transport: &MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
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
						"nova_consulta_serasa_string_retorno": null,
						"logs": [
							{
								"analise_descricao": "Análise CNPJs",
								"regra_descricao": "Cadastro 13 - Análise Sócios - 03",
								"regra_condicao": " SÓCIOS.Homônimo com mesma data de nascimento IGUAL true OU SÓCIOS.Homônimo com mesmo nome da mãe IGUAL true OU SÓCIOS PARTICIPAÇÃO.Tem sanção IGUAL true ",
								"erro": false,
								"alerta": false,
								"liberado": false,
								"transferido": false,
								"erroConsulta": null,
								"bloqueio": false
							},
							{
								"analise_descricao": "Análise CNPJs",
								"regra_descricao": "Cadastro 14 - Análise Sócios - 03",
								"regra_condicao": " SÓCIOS.Homônimo com mesma data de nascimento IGUAL true OU SÓCIOS.Homônimo com mesmo nome da mãe IGUAL true OU SÓCIOS PARTICIPAÇÃO.Tem sanção IGUAL true ",
								"erro": true,
								"alerta": false,
								"liberado": false,
								"transferido": false,
								"erroConsulta": null,
								"bloqueio": false
							}
						]
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
