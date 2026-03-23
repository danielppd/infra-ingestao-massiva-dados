package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ponderada/back/handler"
	"github.com/ponderada/back/queue"
)

func main() {
	// Pega a URL do RabbitMQ das variáveis de ambiente pro Docker, ou usa essa padrão caso eu rode local
	urlRabbit := os.Getenv("RABBITMQ_URL")
	if urlRabbit == "" {
		urlRabbit = "amqp://guest:guest@localhost:5672/"
	}

	// Tentando conectar na nossa fila de mensagens
	filaDeMensagens, erro := queue.NewRabbitMQ(urlRabbit)
	if erro != nil {
		log.Fatalf("Não deu para conectar no RabbitMQ: %v", erro)
	}
	defer filaDeMensagens.Close() // Não podemos esquecer de fechar a conexão no final!

	// Usando o Gin porque vi num tutorial que é ótimo para APIs em Go
	roteador := gin.Default()

	// Cria o nosso controlador passando a fila que acabamos de conectar
	controladorTelemetria := handler.NewTelemetryHandler(filaDeMensagens)
	
	// Rota do tipo POST para receber os dados do sensor
	roteador.POST("/telemetry", controladorTelemetria.PostTelemetry)

	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	log.Printf("Iniciando o nosso servidor backend na porta %s! Tudo pronto pra receber dados.", porta)
	
	// Sobe o servidor!
	if erro := roteador.Run(":" + porta); erro != nil {
		log.Fatalf("Puts deu erro ao subir o servidor: %v", erro)
	}
}