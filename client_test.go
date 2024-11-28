package vadu_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/contbank/vadu-sdk"
	"github.com/contbank/vadu-sdk/mock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VaduClientTestSuite estrutura do teste
type VaduClientTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	ctx        context.Context
	session    *vadu.Session
	vaduClient *vadu.VaduClient
	logger     *logrus.Logger
}

// Função para rodar os testes
func TestVaduClientTestSuite(t *testing.T) {
	suite.Run(t, new(VaduClientTestSuite))
}

// SetupTest configura o ambiente para os testes
func (s *VaduClientTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	// Mock da configuração de sessão
	config := vadu.Config{
		ClientToken: vadu.String("Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJWYWR1IiwidXNyIjoxMTM5NiwiZW1sIjoiY29uZmlnQGNvbnRiYW5rLmNvbS5iciIsImVtcCI6NjY1MTM2MDl9.U9DbXl6UbNPtv9_ZZjgdgodF-ISQIz_B1NPG0me7b_c"),
		Cookie:      vadu.String("mock-cookie-value"),
	}
	// Configurar o logger
	s.logger = logrus.New()
	s.logger.SetFormatter(&logrus.JSONFormatter{}) // Formato estruturado JSON
	s.logger.SetOutput(ioutil.Discard)             // Descartar logs durante testes (opcional)

	// Criar a instância de Session com a configuração fornecida
	session, err := vadu.NewSession(config)
	s.assert.NoError(err) // Verificar se a criação da sessão foi bem-sucedida

	// Atribuir a sessão ao campo s.session
	s.session = session
}

// TestListaGruposAnalise Teste para a função ListaGruposAnalise
func (s *VaduClientTestSuite) TestListaGruposAnalise() {
	// Usando o mock de ListaGruposAnalise
	httpClient := mock.ListaGruposAnaliseMock()

	// Criar o cliente Vadu com o mock HTTP
	s.vaduClient = vadu.NewVaduClient(httpClient, *s.session, s.logger)

	// Chama a função ListaGruposAnalise
	grupos, err := s.vaduClient.ListaGruposAnalise(s.ctx)

	// Verificar se não ocorreu erro e se a resposta está correta
	s.assert.NoError(err)
	s.assert.Len(grupos, 1)
	s.assert.Equal(10802, grupos[0].IDGrupoAnalise)
	s.assert.Equal("Análises CNPJs", grupos[0].NomeGrupoAnalise)
}

// TestListaGruposAnaliseErro Teste para quando a API retorna um erro (status não OK)
func (s *VaduClientTestSuite) TestListaGruposAnaliseErro() {
	// Configurar transporte mockado para simular erro na API
	mockTransport := &http.Client{
		Transport: &mock.MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Simulando um erro 500 (Internal Server Error)
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(strings.NewReader(`{"erro":"Erro interno no servidor"}`)),
				}, nil
			},
		},
	}

	// Substituir cliente HTTP por nosso mock
	s.vaduClient = vadu.NewVaduClient(mockTransport, *s.session, s.logger)

	// Chama a função ListaGruposAnalise
	grupos, err := s.vaduClient.ListaGruposAnalise(s.ctx)

	// Verificar se ocorreu erro
	s.assert.Error(err)
	s.assert.Nil(grupos)
	s.assert.Contains(err.Error(), "erro ao listar grupos de análise")
}

// TestEnviaCNPJsParaAnalise Teste para a função EnviaCNPJsParaAnalise
func (s *VaduClientTestSuite) TestEnviaCNPJsParaAnalise() {
	// Usando o mock de EnviaCNPJsParaAnalise
	httpClient := mock.EnviaCNPJsParaAnaliseMock()

	// Criar o cliente Vadu com o mock HTTP
	s.vaduClient = vadu.NewVaduClient(httpClient, *s.session, s.logger)

	// Definir dados para análise
	postBack := &vadu.PostBack{
		URL:              "https://webhook.site/4c8e0993-5d9c-4f3e-b1c5-bff567be9e78",
		Token:            "Bearer RVkwZWMtITM6NzBrNzhQODJGIyUrrewdRERESdfsaRERER5MDN0Ulo2MkM0LTZhOWU1clFCVSQ2ZEEhJUQtIQ==",
		TipoDadosRetorno: 1,
	}

	// Chama a função EnviaCNPJsParaAnalise
	response, err := s.vaduClient.EnviaCNPJsParaAnalise("33011770000199", 10802, []string{"98960887000164"}, postBack)

	// Verificar se não ocorreu erro e se a resposta está correta
	s.assert.NoError(err)
	s.assert.Equal(4768906, response.AnaliseID)
	s.assert.Equal(1, response.QuantidadeCNPJ)
	s.assert.Equal(0, response.QuantidadeCPF)
	s.assert.Equal("Contbank - Usuário p/Integração Não excluir", response.Usuario)
}

// TestEnviaCNPJsParaAnaliseErro Teste para quando a API retorna um erro (status não OK) na função EnviaCNPJsParaAnalise
func (s *VaduClientTestSuite) TestEnviaCNPJsParaAnaliseErro() {
	// Configurar transporte mockado para simular erro na API
	mockTransport := &http.Client{
		Transport: &mock.MockAuthHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// Simulando um erro 500 (Internal Server Error)
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(strings.NewReader(`{"erro":"Erro interno no servidor"}`)),
				}, nil
			},
		},
	}

	// Substituir cliente HTTP por nosso mock
	s.vaduClient = vadu.NewVaduClient(mockTransport, *s.session, s.logger)

	// Chama a função EnviaCNPJsParaAnalise
	response, err := s.vaduClient.EnviaCNPJsParaAnalise("33011770000199", 10802, []string{"98960887000164"}, nil)

	// Verificar se ocorreu erro
	s.assert.Error(err)
	s.assert.Nil(response)
	s.assert.Contains(err.Error(), "erro ao enviar CNPJs para análise")
}

// TestEnviaCNPJsComDadosParaAnalise
func (s *VaduClientTestSuite) TestEnviaCNPJsComDadosParaAnalise() {
	// Usando o mock de EnviaCNPJsComDadosParaAnalise
	httpClient := mock.EnviaCNPJsParaAnaliseMock()

	// Criar o cliente Vadu com o mock HTTP
	s.vaduClient = vadu.NewVaduClient(httpClient, *s.session, s.logger)

	// Definir os dados para envio (considerando o layout correto)
	listaDados := []vadu.DadosIntegracao{
		{
			CNPJCPF:                       "98960887000164",
			AtivoTotal:                    824167.11,
			AtivoCirculante:               575997.67,
			AtivoNaoCirculante:            248167.44,
			AtivoRealizavelLongoPrazo:     248167.44,
			DeducaoReceitaBruta:           744409.08,
			DepreciacaoBens:               0,
			Despesas:                      744409.08,
			DisponivelCaixa:               470525.47,
			Emprestimo:                    0,
			EstoqueBalanco:                0,
			LucroLiquido:                  2953887.03,
			PassivoCirculante:             771911.23,
			PassivoNaoCirculante:          0,
			PassivoTotal:                  771911.23,
			PatrimonioLiquido:             52255.88,
			ReceitaLiquida:                4610049.5,
			ReceitaBruta:                  5054759.41,
			VendasLiquidas:                4610049.5,
			ScoreExterno:                  761,
			ProbabilidadeInadimplencia:    14,
			DividasBaixasPrejuizo:         0,
			QuantidadeInstituicoes:        0,
			LimiteCreditoVencimentoAte360: 0,
			CreditosVencerAte30Dias:       0,
			Falencia:                      0,
			ChequeSemFundos:               0,
			FaturamentoMedioMensal:        384170.7917,
			CapitalGiroSCR:                0,
			CapitalGiroLiquido:            -195913.56,
			CapitalGiroProprio:            -771909.23,
			NecessidadeCapitalGiro:        0.16147529,
			LiquidezCorrente:              0.746196774,
			LiquidezSeca:                  0.746196774,
			LiquidezGeral:                 1.067696748,
			LiquidezImediata:              0.609559042,
			GrauSolvencia:                 3.826718559,
			Endividamento:                 0.936595529,
			DependenciaRecursosTerceiros:  14.77175832,
			EndividamentoCurtoPrazo:       0.936595529,
			NivelImobilizacao:             4.749081634,
			GrauDependenciaBancaria:       0,
			RetornoPatrimonioLiquidoROE:   56.52736171,
			GiroAtivo:                     67.12303043,
			RetornoSobreAtivoRAO:          3.584087492,
			RetornoSobreVendas:            0.640749526,
			MargemOperacional:             4310350.33,
			RatingExterno:                 "A",
		},
	}

	// Chama a função EnviaCNPJsComDadosParaAnalise
	response, err := s.vaduClient.EnviaCNPJsComDadosParaAnalise("33011770000199", 10802, listaDados)

	// Verificar se não ocorreu erro e se a resposta está correta
	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(4768906, response.AnaliseID)
	s.assert.Equal(1, response.QuantidadeCNPJ)
	s.assert.Equal(0, response.QuantidadeCPF)
	s.assert.Equal("Contbank - Usuário p/Integração Não excluir", response.Usuario)
}

// TestPegaStatusAnalise
func (s *VaduClientTestSuite) TestPegaStatusAnalise() {
	// Usando o mock de PegaStatusAnalise
	httpClient := mock.PegaStatusAnaliseMock()

	s.vaduClient = vadu.NewVaduClient(httpClient, *s.session, s.logger)

	// Chama a função PegaStatusAnalise
	status, err := s.vaduClient.PegaStatusAnalise(4768906)

	// Verifica se não ocorreu erro
	s.assert.NoError(err)

	// Verifica se a resposta está correta
	s.assert.NotNil(status)
	s.assert.Equal(1, status.QuantidadeCNPJsCPFs)
	s.assert.Equal(0, status.QuantidadeConsultasReceita)
	s.assert.Equal(100, status.PercentualConcluido)
	s.assert.False(status.FinalizandoArquivo)
	s.assert.True(status.Concluido)
}

// TestPegaResumoAnalise
func (s *VaduClientTestSuite) TestPegaResumoAnalise() {
	// Usando o mock de PegaResumoAnalise
	httpClient := mock.PegaResumoAnaliseMock()

	// Criar o cliente Vadu com o mock HTTP
	s.vaduClient = vadu.NewVaduClient(httpClient, *s.session, s.logger)

	// Chama a função PegaResumoAnalise
	resumo, err := s.vaduClient.PegaResumoAnalise(4768906)

	// Verifica se não ocorreu erro
	s.assert.NoError(err)

	// Verifica se a resposta está correta
	s.assert.NotNil(resumo)
	s.assert.Equal(4768906, resumo.AnaliseID)
	s.assert.Equal(1, resumo.QuantidadeCNPJ)
	s.assert.Equal("33011770000199", resumo.CNPJEmpresa)
	s.assert.True(resumo.Concluido)
	s.assert.False(resumo.Erro)
	s.assert.True(resumo.Alerta)
	s.assert.Equal(500, resumo.RatingValor)
	s.assert.Equal("B (500)", resumo.RatingSigla)
}

// TestListaResumoCNPJs
func (s *VaduClientTestSuite) TestListaResumoCNPJs() {
	// Usando o mock de ListaResumoCNPJs
	httpClient := mock.ListaResumoCNPJsMock()

	// Criar o cliente Vadu com o mock HTTP
	s.vaduClient = vadu.NewVaduClient(httpClient, *s.session, s.logger)

	// Chama a função ListaResumoCNPJs
	resumos, err := s.vaduClient.ListaResumoCNPJs(4768906)

	// Verifica se não ocorreu erro
	s.assert.NoError(err)

	// Verifica se a resposta está correta
	s.assert.NotNil(resumos)
	s.assert.Len(resumos, 1) // Verifica que existe um resumo
	s.assert.Equal(4768906, resumos[0].AnaliseID)
	s.assert.Equal(28883956, resumos[0].AnaliseCNPJCPFID)
	s.assert.Equal("98960887000164", resumos[0].CNPJCPF)
	s.assert.Equal("WEBSOLUTIONS LTDA", resumos[0].Nome)
	s.assert.False(resumos[0].Erro)
	s.assert.True(resumos[0].Alerta)
	s.assert.False(resumos[0].Bloqueio)
	s.assert.Equal(500, resumos[0].Rating)
	s.assert.Equal("B (500)", resumos[0].RatingSigla)
	s.assert.Equal("CCB--> 40% do faturamento mensal || ANT--> 20% do faturamento mensal", resumos[0].RatingDescricao)
	s.assert.False(resumos[0].NovaConsultaSerasa)
	s.assert.Equal("", resumos[0].NovaConsultaSerasaString)
	s.assert.Equal("Consulta Serasa", resumos[0].OrigemConsultaSerasaTexto)
}

// TestListaResumoCNPJsDetalhado
func (s *VaduClientTestSuite) TestListaResumoCNPJsDetalhado() {
	// Usando o mock
	httpClient := mock.ListaResumoCNPJsDetalhadoMock()

	// Criar o cliente Vadu com o mock HTTP
	s.vaduClient = vadu.NewVaduClient(httpClient, *s.session, s.logger)

	// Chama a função ListaResumoCNPJsDetalhado
	resumos, err := s.vaduClient.ListaResumoCNPJsDetalhado(4768906)

	// Verifica se não ocorreu erro
	s.assert.NoError(err)

	// Verifica se a resposta está correta
	s.assert.NotNil(resumos)
	s.assert.Len(resumos, 1)
	s.assert.Equal(4768906, resumos[0].AnaliseID)
	s.assert.Equal("37697591000108", resumos[0].CNPJCPF)
	s.assert.Equal("WEBTECH SOLUTIONS CONSULTORIA E INFORMATICA LTDA", resumos[0].Nome)
	s.assert.False(resumos[0].Erro)
	s.assert.True(resumos[0].Alerta)
	s.assert.Equal(500, resumos[0].Rating)
	s.assert.Equal("B (500)", resumos[0].RatingSigla)
	s.assert.Len(resumos[0].Logs, 1)
	s.assert.Equal("Análise CNPJs", resumos[0].Logs[0].AnaliseDescricao)
	s.assert.Equal("Cadastro 14 - Análise Sócios - 03", resumos[0].Logs[0].RegraDescricao)
}
