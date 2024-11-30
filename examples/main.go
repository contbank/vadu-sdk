package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/contbank/vadu-sdk"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configurar logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Configuração da sessão
	config := vadu.Config{
		APIEndpoint:   vadu.String("https://www.vadu.com.br"),
		ClientToken:   vadu.String("chave api"),
		Cookie:        vadu.String("vadu-sdk"),
		LoginEndpoint: vadu.String("https://www.vadu.com.br/vadu.dll/Autenticacao/JSONPegarToken"),
	}

	// Criar sessão
	session, err := vadu.NewSession(config)
	if err != nil {
		logger.WithError(err).Fatal("Erro ao criar sessão")
	}

	// Configurar cliente HTTP
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	logger.WithFields(logrus.Fields{
		"httpClient": httpClient,
	}).Info("HTTP Client configurado com sucesso")

	// Criar a instância de autenticação
	auth := vadu.NewAuthentication(httpClient, *session, logger)

	// Criar cliente Vadu
	vaduClient := vadu.NewVaduClient(httpClient, *session, logger)

	// Contexto para a chamada
	ctx := context.Background()

	// Testar autenticação
	token, err := auth.Token(ctx)
	if err != nil {
		logger.WithError(err).Fatal("Erro ao autenticar")
	}
	logger.WithField("token", token[:5]+"...").Info("Token autenticado com sucesso")

	// Testar a função ListaGruposAnalise
	grupos, err := vaduClient.ListaGruposAnalise(ctx, auth)
	if err != nil {
		logger.WithError(err).Fatal("Erro ao listar grupos de análise")
	}

	// Imprimir os grupos de análise recebidos
	for _, grupo := range grupos {
		fmt.Printf("ID: %d, Nome: %s, Rating Mínimo: %d, Rating Máximo: %d\n",
			grupo.IDGrupoAnalise, grupo.NomeGrupoAnalise, grupo.RatingMinimo, grupo.RatingMaximo)
	}

	log.Println("Teste concluído com sucesso")
}
