package vadu

import (
	"time"
)

// GrupoAnalise representa a estrutura de um grupo de análise.
type GrupoAnalise struct {
	IDGrupoAnalise       int    `json:"id_grupo_analise"`
	NomeGrupoAnalise     string `json:"nome_grupo_analise"`
	RatingStart          int    `json:"rating_start"`
	RatingMinimo         int    `json:"rating_minimo"`
	RatingMaximo         int    `json:"rating_maximo"`
	QuantidadeAnalises   int    `json:"quantidade_analises"`
	QuantidadeRegras     int    `json:"quantidade_regras"`
	QuantidadeValidacoes int    `json:"quantidade_validacoes"`
}

// PostBack define a estrutura opcional para o campo de postback
type PostBack struct {
	URL              string `json:"url"`
	Token            string `json:"token"`
	TipoDadosRetorno int    `json:"tipoDadosRetorno"`
}

// EnviaCNPJsRequest define os dados que são enviados na requisição para a análise de CNPJs
type EnviaCNPJsRequest struct {
	CNPJEmpresa    string    `json:"cnpj_empresa"`
	IDGrupoAnalise int       `json:"id_grupo_analise"`
	ListaCNPJCPF   []string  `json:"lista_cnpj_cpf"`
	PostBack       *PostBack `json:"postBack,omitempty"` // PostBack é opcional
}

// EnviaCNPJsResponse define a estrutura da resposta da API de envio de CNPJs
type EnviaCNPJsResponse struct {
	AnaliseID        int       `json:"analise_id"`
	QuantidadeCNPJ   int       `json:"quantidade_cnpj"`
	QuantidadeCPF    int       `json:"quantidade_cpf"`
	Usuario          string    `json:"usuario"`
	DataHoraEnvio    time.Time `json:"data_hora_envio"`
	IDGrupoAnalise   int       `json:"id_grupo_analise"`
	NomeLote         string    `json:"nome_lote"`
	NomeGrupoAnalise string    `json:"nome_grupo_analise"`
}

// DadosIntegracao representa os dados detalhados enviados para análise.
type DadosIntegracao struct {
	CNPJCPF                       string  `json:"cnpjcpf"`
	AtivoTotal                    float64 `json:"ativoTotal"`
	AtivoCirculante               float64 `json:"ativoCirculante"`
	AtivoNaoCirculante            float64 `json:"ativoNaoCirculante"`
	AtivoRealizavelLongoPrazo     float64 `json:"ativoRealizavelLongoPrazo"`
	DeducaoReceitaBruta           float64 `json:"deducaoReceitaBruta"`
	DepreciacaoBens               float64 `json:"depreciacaoBens"`
	Despesas                      float64 `json:"despesas"`
	DisponivelCaixa               float64 `json:"disponivelCaixa"`
	Emprestimo                    float64 `json:"emprestimo"`
	EstoqueBalanco                float64 `json:"estoqueBalanco"`
	LucroLiquido                  float64 `json:"lucroLiquido"`
	PassivoCirculante             float64 `json:"passivoCirculante"`
	PassivoNaoCirculante          float64 `json:"passivoNaoCirculante"`
	PassivoTotal                  float64 `json:"passivoTotal"`
	PatrimonioLiquido             float64 `json:"patrimonioLiquido"`
	ReceitaLiquida                float64 `json:"receitaLiquida"`
	ReceitaBruta                  float64 `json:"receitaBruta"`
	VendasLiquidas                float64 `json:"vendasLiquidas"`
	ScoreExterno                  int     `json:"scoreExterno"`
	ProbabilidadeInadimplencia    int     `json:"probabilidadeInadimplencia"`
	DividasBaixasPrejuizo         float64 `json:"dividasBaixasPrejuizo"`
	QuantidadeInstituicoes        int     `json:"quantidadeInstituicoes"`
	LimiteCreditoVencimentoAte360 float64 `json:"limiteCreditoVencimentoAte360Dias"`
	CreditosVencerAte30Dias       float64 `json:"creditosVencerAte30Dias"`
	Falencia                      int     `json:"falencia"`
	ChequeSemFundos               int     `json:"chequeSemFundos"`
	FaturamentoMedioMensal        float64 `json:"faturamentoMedioMensal"`
	CapitalGiroSCR                float64 `json:"capitalGiroSCR"`
	CapitalGiroLiquido            float64 `json:"capitalGiroLiquido"`
	CapitalGiroProprio            float64 `json:"capitalGiroProprio"`
	NecessidadeCapitalGiro        float64 `json:"necessidadeCapitalGiro"`
	LiquidezCorrente              float64 `json:"liquidezCorrente"`
	LiquidezSeca                  float64 `json:"liquidezSeca"`
	LiquidezGeral                 float64 `json:"liquidezGeral"`
	LiquidezImediata              float64 `json:"liquidezImediata"`
	GrauSolvencia                 float64 `json:"grauSolvencia"`
	Endividamento                 float64 `json:"endividamento"`
	DependenciaRecursosTerceiros  float64 `json:"dependeciaRecursosTerceiros"`
	EndividamentoCurtoPrazo       float64 `json:"endividamentoCurtoPrazo"`
	NivelImobilizacao             float64 `json:"nivelImobilizacao"`
	GrauDependenciaBancaria       float64 `json:"grauDependenciaBancaria"`
	RetornoPatrimonioLiquidoROE   float64 `json:"retornoPatrimonioLiquidoROE"`
	GiroAtivo                     float64 `json:"giroAtivo"`
	RetornoSobreAtivoRAO          float64 `json:"retornoSobreAtivoRAO"`
	RetornoSobreVendas            float64 `json:"retornoSobreVendas"`
	MargemOperacional             float64 `json:"margemOperacional"`
	RatingExterno                 string  `json:"ratingExterno"`
}

type LogAnalise struct {
	AnaliseDescricao string `json:"analise_descricao"`
	RegraDescricao   string `json:"regra_descricao"`
	RegraCondicao    string `json:"regra_condicao"`
	Erro             bool   `json:"erro"`
	Alerta           bool   `json:"alerta"`
	Liberado         bool   `json:"liberado"`
	Transferido      bool   `json:"transferido"`
	ErroConsulta     *bool  `json:"erroConsulta"`
	Bloqueio         bool   `json:"bloqueio"`
}

type ResumoCNPJDatalhado struct {
	AnaliseID                 int          `json:"analise_id"`
	AnaliseCNPJCPFID          int          `json:"analise_cnpj_cpf_id"`
	CNPJCPF                   string       `json:"cnpj_cpf"`
	Nome                      string       `json:"nome"`
	Erro                      bool         `json:"erro"`
	Alerta                    bool         `json:"alerta"`
	Bloqueio                  bool         `json:"bloqueio"`
	Rating                    int          `json:"rating"`
	RatingSigla               string       `json:"rating_sigla"`
	RatingDescricao           string       `json:"rating_descricao"`
	Rating2                   int          `json:"rating2"`
	Rating2Sigla              string       `json:"rating2_sigla"`
	Rating2Descricao          string       `json:"rating2_descricao"`
	NovaConsultaSerasa        bool         `json:"nova_consulta_serasa"`
	NovaConsultaSerasaString  *string      `json:"nova_consulta_serasa_string_retorno"`
	FlowSolicitacaoID         int          `json:"flowSolicitacao_id"`
	FlowTarefaNome            string       `json:"flowTarefaNome"`
	FlowTarefaID              int          `json:"flowTarefa_id"`
	OrigemConsultaSerasa      int          `json:"origem_consulta_serasa"`
	OrigemConsultaSerasaTexto string       `json:"origem_consulta_serasa_texto"`
	Logs                      []LogAnalise `json:"logs"`
}

// EnviaCNPJsComDadosRequest representa a estrutura do corpo da requisição para envio de CNPJs com dados detalhados.
type EnviaCNPJsComDadosRequest struct {
	CNPJEmpresa                 string            `json:"cnpj_empresa"`
	IDGrupoAnalise              int               `json:"id_grupo_analise"`
	ListaCNPJCPFDadosIntegracao []DadosIntegracao `json:"lista_cnpj_cpf_dados_integracao"`
	PostBack                    *PostBack         `json:"postBack,omitempty"` // PostBack é opcional
}

// StatusAnalise representa a resposta da API de status de análise
type StatusAnalise struct {
	QuantidadeCNPJsCPFs           int  `json:"quantidade_cnpj_cpf"`
	QuantidadeConsultasReceita    int  `json:"quantidade_consultas_receita"`
	PercentualConsultasReceita    int  `json:"percentual_consultas_receita"`
	QuantidadeCNPJsCPFsConcluidos int  `json:"quantidade_cnpj_cpf_concluidos"`
	PercentualConcluido           int  `json:"percentual_concluido"`
	FinalizandoArquivo            bool `json:"finalizando_arquivo"`
	Concluido                     bool `json:"concluido"`
}

// ResumoAnalise representa a resposta da API de resumo da análise
type ResumoAnalise struct {
	AnaliseID              int    `json:"analise_id"`
	QuantidadeCNPJ         int    `json:"quantidade_cnpj"`
	QuantidadeCPF          int    `json:"quantidade_cpf"`
	CNPJEmpresa            string `json:"cnpj_empresa"`
	Usuario                string `json:"usuario"`
	DataHoraEnvio          string `json:"data_hora_envio"`
	DataHoraConclusao      string `json:"data_hora_conclusao"`
	Concluido              bool   `json:"concluido"`
	Erro                   bool   `json:"erro"`
	Alerta                 bool   `json:"alerta"`
	Bloqueio               bool   `json:"bloqueio"`
	QuantidadeCNPJAlerta   int    `json:"quantidade_cnpj_alerta"`
	QuantidadeCNPJBloqueio int    `json:"quantidade_cnpj_bloqueio"`
	QuantidadeCPFAlerta    int    `json:"quantidade_cpf_alerta"`
	QuantidadeCPFBloqueio  int    `json:"quantidade_cpf_bloqueio"`
	IDGrupoAnalise         int    `json:"id_grupo_analise"`
	NomeGrupoAnalise       string `json:"nome_grupo_analise"`
	RatingValor            int    `json:"rating_valor"`
	RatingSigla            string `json:"rating_sigla"`
	RatingDescricao        string `json:"rating_descricao"`
	Rating2Valor           int    `json:"rating2_valor"`
	Rating2Sigla           string `json:"rating2_sigla"`
	Rating2Descricao       string `json:"rating2_descricao"`
	NomeLote               string `json:"nome_lote"`
}

// Estrutura para armazenar os dados do resumo de CNPJs analisados
type ResumoCNPJ struct {
	AnaliseID                 int    `json:"analise_id"`
	AnaliseCNPJCPFID          int    `json:"analise_cnpj_cpf_id"`
	CNPJCPF                   string `json:"cnpj_cpf"`
	Nome                      string `json:"nome"`
	Erro                      bool   `json:"erro"`
	Alerta                    bool   `json:"alerta"`
	Bloqueio                  bool   `json:"bloqueio"`
	Rating                    int    `json:"rating"`
	RatingSigla               string `json:"rating_sigla"`
	RatingDescricao           string `json:"rating_descricao"`
	Rating2                   int    `json:"rating2"`
	Rating2Sigla              string `json:"rating2_sigla"`
	Rating2Descricao          string `json:"rating2_descricao"`
	NovaConsultaSerasa        bool   `json:"nova_consulta_serasa"`
	NovaConsultaSerasaString  string `json:"nova_consulta_serasa_string_retorno"`
	FlowSolicitacaoID         int    `json:"flowSolicitacao_id"`
	FlowTarefaNome            string `json:"flowTarefaNome"`
	FlowTarefaID              int    `json:"flowTarefa_id"`
	OrigemConsultaSerasa      int    `json:"origem_consulta_serasa"`
	OrigemConsultaSerasaTexto string `json:"origem_consulta_serasa_texto"`
}
